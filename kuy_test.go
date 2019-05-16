package kuy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_kuy(t *testing.T) {

	var (
		numOfUser         = 500
		maxItem           = 5
		expectedNumOfPool = numOfUser / maxItem
	)

	e := New(maxItem)

	for i := 0; i < numOfUser; i++ {
		e.Join(i)
	}

	nop := e.GetNumberOfPools()
	require.Equal(t, expectedNumOfPool, nop)

}
