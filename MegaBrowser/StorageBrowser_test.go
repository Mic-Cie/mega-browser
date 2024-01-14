package megabrowser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t3rm1n4l/go-mega"
)

var errLogin = fmt.Errorf("mock login error")

type mockClient struct {
	errLogin error
}

func TestShouldFailCreatingBrowserIfFailedToLogin(t *testing.T) {
	mockClient := mockClient{
		errLogin: errLogin,
	}
	storageBrowser, err := NewMegaBrowser(&mockClient)
	assert.Nil(t, storageBrowser)
	assert.Equal(t, errLogin, err)
}

func TestShouldGetObjectNode(t *testing.T) {
	storageBrowser, errBrowserCreate := NewMegaBrowser(&mockClient{})
	require.Nil(t, errBrowserCreate)
	node, err := storageBrowser.GetObjectNode("")
	assert.Equal(t, "", node)
	assert.Nil(t, err)
}

func (m *mockClient) Login(login string, pass string) error {
	return m.errLogin
}

func (m *mockClient) GetFS() *mega.MegaFS {
	return nil
}
