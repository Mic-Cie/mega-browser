package megabrowser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t3rm1n4l/go-mega"
)

func TestShouldConvertStructArrayToInterfaceArray(t *testing.T) {
	nodes := []*mega.Node{
		{}, {},
	}

	convertedNodes := nodeStructArrToInterfaceArr(nodes)

	assert.Equal(t, len(nodes), len(convertedNodes))
}
