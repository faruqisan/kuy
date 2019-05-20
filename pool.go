package kuy

import (
	"sync"
)

type (
	pool struct {
		id      string
		maxItem int

		m     sync.Mutex
		items []interface{}

		expireWaitCount int

		respChan chan PoolResp
	}

	// PoolResp function define response when item joining the pool
	PoolResp struct {
		PoolID   string
		IsFull   bool
		TimeIsUp bool
		Items    []interface{}
	}
)

// NewPool func create new pool
func newPool(id string, maxItem int) *pool {
	return &pool{
		id:       id,
		maxItem:  maxItem,
		respChan: make(chan PoolResp, maxItem),
	}
}

func (p *pool) add(item interface{}) chan PoolResp {
	p.m.Lock()
	defer func() {
		if len(p.items) == p.maxItem {
			p.respChan <- PoolResp{
				PoolID: p.id,
				IsFull: true,
				Items:  p.items,
			}
		}
		p.m.Unlock()
	}()
	p.items = append(p.items, item)

	return p.respChan

}

func (p *pool) ableToJoin() bool {
	p.m.Lock()
	defer p.m.Unlock()

	return len(p.items) < p.maxItem
}
