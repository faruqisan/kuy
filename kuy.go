package kuy

import (
	"log"

	"github.com/gofrs/uuid"
)

// Engine struct hold required data, and act as function receiver
type Engine struct {
	maxItem int
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
	)

	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}

	if len(e.pools) == 0 {
		p = newPool(id.String(), e.maxItem)
		e.pools = append(e.pools, p)
		return p
	}

	for _, v := range e.pools {
		if v.ableToJoin() {
			return v
		}
	}

	if !joined {
		p = newPool(id.String(), e.maxItem)
		e.pools = append(e.pools, p)
		return p
	}

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
	return len(e.pools)
}
