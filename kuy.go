package kuy

import (
	"log"
	"sync"

	"github.com/gofrs/uuid"
)

// Engine struct hold required data, and act as function receiver
type Engine struct {
	maxItem int
	mutex   sync.Mutex
	pools   []*pool
}

// New function return engine struct
func New(maxItem int) *Engine {
	return &Engine{
		maxItem: maxItem,
	}
}

func (e *Engine) getAvailablePool() *pool {
	var (
		p      *pool
		joined bool
		nop    = e.GetNumberOfPools()
	)

	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}

	if nop == 0 {
		return e.createPool(id.String())
	}

	for _, v := range e.pools {
		if v.ableToJoin() {
			return v
		}
	}

	if !joined {
		return e.createPool(id.String())
	}

	return p
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
	p := e.getAvailablePool()
	return p.add(item)
}

// GetNumberOfPools return number of pools
func (e *Engine) GetNumberOfPools() int {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return len(e.pools)
}
