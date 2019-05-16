package kuy

import "sync"

type (
	pool struct {
		id      string
		maxItem int
		m       sync.Mutex
		items   []interface{}

		fullChan chan PoolResp
	}

	// PoolResp function define response when item joining the pool
	PoolResp struct {
		IsFull bool
		Items  []interface{}
	}
)

// NewPool func create new pool
func newPool(id string, maxItem int) *pool {
	return &pool{
		maxItem:  maxItem,
		fullChan: make(chan PoolResp, 1),
	}
}

func (p *pool) add(item interface{}) chan PoolResp {
	p.m.Lock()
	defer func() {
		if len(p.items) == p.maxItem {
			p.fullChan <- PoolResp{
				IsFull: true,
				Items:  p.items,
			}
		}
		p.m.Unlock()
	}()

	p.items = append(p.items, item)

	return p.fullChan

}

func (p *pool) ableToJoin() bool {
	p.m.Lock()
	defer p.m.Unlock()

	return len(p.items) < p.maxItem
}
