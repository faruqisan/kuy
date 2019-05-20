package kuy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_kuy(t *testing.T) {

	var (
		numOfUser         = 20
		maxItem           = 10
		expectedNumOfPool = numOfUser / maxItem
	)

	e := New(Option{
		MaxItem: maxItem,
	})

	for i := 0; i < numOfUser; i++ {
		e.Join(i)
	}

	nop := e.GetNumberOfPools()
	require.Equal(t, expectedNumOfPool, nop)

}
