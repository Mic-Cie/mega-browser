package megabrowser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t3rm1n4l/go-mega"
)

type mockNode struct {
	name     string
	nodeType int
	hash     string
}

func TestShouldConvertStructArrayToInterfaceArray(t *testing.T) {
	nodes := []*mega.Node{
		{}, {},
	}

	convertedNodes := nodeStructArrToInterfaceArr(nodes)

	assert.Equal(t, len(nodes), len(convertedNodes))
}

func TestShouldReturnNodeHashIfExistsOnTheList(t *testing.T) {
	expName := "expected file name"
	expHash := "exphash"
	nodes := []Node{
		&mockNode{
			name:     "other name",
			nodeType: fileType,
			hash:     "hash",
		},
		&mockNode{
			name:     expName,
			nodeType: fileType,
			hash:     expHash,
		},
	}

	hash, err := getNodeHashOfExpectedFile(expName, &nodes)

	assert.Equal(t, expHash, hash)
	assert.Nil(t, err)
}

func TestShouldFailIfFileNotExistOnTheList(t *testing.T) {
	expName := "expected file name"
	nodes := []Node{
		&mockNode{
			name:     "other name",
			nodeType: fileType,
			hash:     "hash",
		},
	}

	hash, err := getNodeHashOfExpectedFile(expName, &nodes)

	assert.Empty(t, hash)
	assert.Equal(t, fmt.Errorf("could not find file: %s", expName), err)
}

func (m *mockNode) GetName() string {
	return m.name
}

func (m *mockNode) GetType() int {
	return m.nodeType
}

func (m *mockNode) GetHash() string {
	return m.hash
}
