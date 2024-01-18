package megabrowser

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t3rm1n4l/go-mega"
)

const (
	expRootNodeHash = "rootnodehash"
	expFileHash     = "expfilehash"
	expDirHash      = "expdirhash"
	expFileName     = "expFile"
	expDirName      = "expDir"
)

var (
	errLogin       = fmt.Errorf("mock login error")
	errGetChildren = fmt.Errorf("mock get children")
)

type mockClient struct {
	errLogin error
	fs       Fs
}

type mockFs struct {
	children       []*mega.Node
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
		errGetChildren: nil,
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

func TestShouldSuccessfullyGetObjectNode(t *testing.T) {
	mockFs := mockFs{
		errGetChildren: nil,
	}
	mockClient := mockClient{
		nil,
		&mockFs,
	}
	storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
	storageBrowser.getRootNodeHash = getRootNodeHash
	storageBrowser.getChildren = mockGetChildren

	result, err := storageBrowser.GetObjectNode(filepath.Join(expDirName, expFileName))

	assert.Equal(t, expFileHash, result)
	assert.Nil(t, err)
}

func TestShouldFailGettingObjectNodeWhenCouldNotFindExpectedDirectory(t *testing.T) {
	mockFs := mockFs{
		errGetChildren: nil,
	}
	mockClient := mockClient{
		nil,
		&mockFs,
	}
	storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
	storageBrowser.getRootNodeHash = getRootNodeHash
	storageBrowser.getChildren = mockGetChildren

	result, err := storageBrowser.GetObjectNode(filepath.Join("unexpectedDir", expFileName))

	assert.Empty(t, result)
	assert.Equal(t, fmt.Errorf("could not find directory: unexpectedDir"), err)
}

func TestShouldFailGettingObjectNodeWhenCouldNotFindExpectedFile(t *testing.T) {
	mockFs := mockFs{
		errGetChildren: nil,
	}
	mockClient := mockClient{
		nil,
		&mockFs,
	}
	storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
	storageBrowser.getRootNodeHash = getRootNodeHash
	storageBrowser.getChildren = mockGetChildren

	result, err := storageBrowser.GetObjectNode(filepath.Join(expDirName, "unexpFile"))

	assert.Empty(t, result)
	assert.Equal(t, fmt.Errorf("could not find file: unexpFile"), err)
}

func TestShouldFailGettingObjectNodeWhenBothPathAndSeparatorAreEmpty(t *testing.T) {
	mockFs := mockFs{
		errGetChildren: nil,
	}
	mockClient := mockClient{
		nil,
		&mockFs,
	}
	storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
	storageBrowser.getRootNodeHash = getRootNodeHash
	storageBrowser.getChildren = mockGetChildren
	storageBrowser.targetSeparator = ""

	result, err := storageBrowser.GetObjectNode("")

	assert.Empty(t, "", result)
	assert.Equal(t, fmt.Errorf("could not find object node for "), err)
}

func TestShouldFailWhenGettingObjectNodeForEmptyPath(t *testing.T) {
	storageBrowser := NewMegaBrowser(nil, nil)

	result, err := storageBrowser.GetObjectNode("")

	assert.Empty(t, result)
	assert.Equal(t, fmt.Errorf("trying to find object node for an empty path"), err)
}

func TestShouldFailGettingObjectNodeWhenCouldNotGetChildren(t *testing.T) {
	mockFs := mockFs{
		errGetChildren: errGetChildren,
	}
	storageBrowser := NewMegaBrowser(nil, &mockFs)

	result, err := storageBrowser.GetObjectNode("path")

	assert.Empty(t, result)
	assert.Equal(t, errGetChildren, err)
}

func TestShouldGetChildren(t *testing.T) {
	nodes := []*mega.Node{
		{}, {},
	}
	mockFs := mockFs{
		children:       nodes,
		errGetChildren: nil,
	}

	children, err := getChildren(&mockFs, "path")

	assert.Equal(t, len(nodes), len(children))
	assert.Nil(t, err)
}

func (m *mockClient) Login(login string, pass string) error {
	return m.errLogin
}

func (m *mockFs) GetChildren(node *mega.Node) ([]*mega.Node, error) {
	return m.children, m.errGetChildren
}

func (m *mockFs) GetRoot() *mega.Node {
	return nil
}

func (m *mockFs) HashLookup(string) *mega.Node {
	return nil
}

func mockGetRootNodeHash(nodes []Node) (string, error) {
	return expRootNodeHash, nil
}

func mockGetChildren(fs Fs, nodeHash string) ([]Node, error) {
	if nodeHash == expDirHash {
		return []Node{
			&mockNode{
				name:     expFileName,
				nodeType: fileType,
				hash:     expFileHash,
			},
		}, nil
	} else {
		return []Node{
			&mockNode{
				name:     expDirName,
				nodeType: directoryType,
				hash:     expDirHash,
			},
		}, nil
	}
}
