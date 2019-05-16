package kuy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_pool(t *testing.T) {

	var (
		id      = "a"
		maxItem = 5
	)

	pool := newPool(id, maxItem)

	for i := 0; i < maxItem+1; i++ {
		if pool.ableToJoin() {
			pool.add(i)
		}
	}

	items := pool.items
	require.Equal(t, maxItem, len(items))

}
