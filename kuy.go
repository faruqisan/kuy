package kuy

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	emptyTime = "0001-01-01 00:00:00 +0000 UTC"
)

var (
	defaultWaitPeriod = time.Second * 2
)

// Engine struct hold required data, and act as function receiver
type Engine struct {
	maxItem     int
	waitPeriod  time.Duration
	mutex       sync.Mutex
	pools       []*pool
	expiredPool map[string]struct{}
}

// Option struct define engine option configuration
type Option struct {
	MaxItem    int
	WaitPeriod time.Duration
}

// New function return engine struct
func New(opt Option) *Engine {

	var (
		wp = defaultWaitPeriod
	)

	if opt.WaitPeriod.String() == emptyTime {
		wp = opt.WaitPeriod
	}

	return &Engine{
		maxItem:     opt.MaxItem,
		waitPeriod:  wp,
		expiredPool: make(map[string]struct{}),
	}
}

func (e *Engine) getAvailablePool() *pool {
	var (
		nop = e.GetNumberOfPools()
		id  = uuid.New().ID()
		sID = fmt.Sprintf("%d", id)
	)

	if nop == 0 {
		return e.createPool(sID)
	}

	// TODO: improve find pools
	// currently we just loop through pools on engine
	for _, v := range e.pools {
		if v.ableToJoin() {
			return v
		}
	}

	return e.createPool(sID)
}

func (e *Engine) createPool(id string) *pool {
	p := newPool(id, e.maxItem)

	e.mutex.Lock()
	e.pools = append(e.pools, p)
	e.mutex.Unlock()

	return p
}

// Join function add given item into available pool, returning chanel of PoolResp
// that will notify when pool is full and ready.
func (e *Engine) Join(item interface{}) chan PoolResp {
	var (
		p     = e.getAvailablePool()
		timer = time.NewTimer(e.waitPeriod)
	)

	go func() {
		select {
		case <-timer.C:

			e.mutex.Lock()

			if p.ableToJoin() {

				p.expireWaitCount++

				if _, ok := e.expiredPool[p.id]; !ok {

					p.respChan <- PoolResp{
						PoolID:   p.id,
						TimeIsUp: true,
						Items:    p.items,
					}
					e.expiredPool[p.id] = struct{}{}
				}

				if p.expireWaitCount == len(p.items) {
					// remove items on pool
					p.items = nil
					// remove pool from expired map
					delete(e.expiredPool, p.id)
				}
			}

			e.mutex.Unlock()
			break
		}
	}()

	return p.add(item)
}

// GetNumberOfPools return number of pools
func (e *Engine) GetNumberOfPools() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return len(e.pools)
}
