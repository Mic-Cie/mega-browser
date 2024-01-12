package megabrowser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldGetNilFS(t *testing.T) {
	client := MegaClient{}
	fs := client.GetFS()

	assert.Nil(t, fs)
}
