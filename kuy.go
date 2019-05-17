package kuy

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Engine struct hold required data, and act as function receiver
type Engine struct {
	maxItem    int
	waitPeriod time.Duration
	mutex      sync.Mutex
	pools      []*pool
}

// Option struct define engine option configuration
type Option struct {
	MaxItem    int
	WaitPeriod time.Duration
}

// New function return engine struct
func New(opt Option) *Engine {
	return &Engine{
		maxItem:    opt.MaxItem,
		waitPeriod: opt.WaitPeriod,
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
			if p.ableToJoin() {
				p.respChan <- PoolResp{
					PoolID:   p.id,
					TimeIsUp: true,
					Items:    p.items,
				}
			}
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
