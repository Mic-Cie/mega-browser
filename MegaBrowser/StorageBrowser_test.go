package megabrowser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t3rm1n4l/go-mega"
)

const expRootNodeHash = "rootnodehash"

var (
	errLogin       = fmt.Errorf("mock login error")
	errGetChildren = fmt.Errorf("mock get children")
)

type mockClient struct {
	errLogin error
	fs       Fs
}

type mockFs struct {
	errGetChildren error
}

func TestShouldSuccessfullyInitializeBrowser(t *testing.T) {
	mockFs := mockFs{}
	mockClient := mockClient{
		errLogin: nil,
		fs:       &mockFs,
	}
	storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
	storageBrowser.getRootNodeHash = mockGetRootNodeHash
	err := storageBrowser.Initialize()
	assert.Equal(t, expRootNodeHash, storageBrowser.rootNodeHash)
	assert.Nil(t, err)
}

func TestShouldFailInitializingBrowserIfFailedToLogin(t *testing.T) {
	mockClient := mockClient{
		errLogin: errLogin,
	}
	storageBrowser := NewMegaBrowser(&mockClient, nil)
	err := storageBrowser.Initialize()
	assert.Equal(t, errLogin, err)
}

func TestShouldFailIfCouldNotGetNodeChildren(t *testing.T) {
	mockFs := mockFs{
		errGetChildren: errGetChildren,
	}
	mockClient := mockClient{
		nil,
		&mockFs,
	}
	storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
	err := storageBrowser.Initialize()
	require.Equal(t, errGetChildren, err)
}

func TestShouldFailIfCouldNotGetRootNodeHash(t *testing.T) {
	mockFs := mockFs{
		nil,
	}
	mockClient := mockClient{
		nil,
		&mockFs,
	}
	storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
	err := storageBrowser.Initialize()
	require.Equal(t, errGetRootNodeHash, err)
}

func TestShouldFailIfNoneOfTheNodesIsRoot(t *testing.T) {
	nodes := []Node{
		&mockNode{}, &mockNode{},
	}

	rootNodeHash, err := getRootNodeHash(nodes)

	assert.Empty(t, rootNodeHash)
	assert.Equal(t, errGetRootNodeHash, err)
}

func TestShouldSuccessfullyGetRootNodeHash(t *testing.T) {
	nodes := []Node{
		&mockNode{
			name:     rootNodeName,
			nodeType: 1,
			hash:     expRootNodeHash,
		},
	}

	rootNodeHash, err := getRootNodeHash(nodes)

	assert.Equal(t, expRootNodeHash, rootNodeHash)
	assert.Nil(t, err)
}

func TestShouldNotFailWhenGettingObjectNode(t *testing.T) {
	storageBrowser := NewMegaBrowser(nil, nil)
	_, err := storageBrowser.GetObjectNode("file")
	assert.Nil(t, err)
}

func (m *mockClient) Login(login string, pass string) error {
	return m.errLogin
}

func (m *mockFs) GetChildren(node *mega.Node) ([]*mega.Node, error) {
	return nil, m.errGetChildren
}

func (m *mockFs) GetRoot() *mega.Node {
	return nil
}

func mockGetRootNodeHash(nodes []Node) (string, error) {
	return expRootNodeHash, nil
}
