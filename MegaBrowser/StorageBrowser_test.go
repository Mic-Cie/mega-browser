package megabrowser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldGetObjectNode(t *testing.T) {
	storageBrowser, _ := NewMegaBrowser()
	node, err := storageBrowser.GetObjectNode("")
	assert.Equal(t, "", node)
	assert.Nil(t, err)
}
