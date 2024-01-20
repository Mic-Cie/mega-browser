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
	errGetChildren = fmt.Errorf("mock get children error")
	errDownload    = fmt.Errorf("mock download error")
)

type mockClient struct {
	errLogin    error
	errDownload error
}

type mockFs struct {
	children       []*mega.Node
	errGetChildren error
}

func TestInitializeStorageBrowser(t *testing.T) {
	tests := []struct {
		name                    string
		loginError              error
		getChildrenError        error
		getRootNodeHashFunction getRootNodeHashFunc
		expRootNodeHash         string
		expErr                  error
	}{
		{
			name:                    "should successfully initialize browser, if all inputs are valid",
			loginError:              nil,
			getChildrenError:        nil,
			getRootNodeHashFunction: mockGetRootNodeHash,
			expRootNodeHash:         expRootNodeHash,
			expErr:                  nil,
		},
		{
			name:                    "should fail, if could not login to Mega",
			loginError:              errLogin,
			getChildrenError:        nil,
			getRootNodeHashFunction: mockGetRootNodeHash,
			expRootNodeHash:         "",
			expErr:                  errLogin,
		},
		{
			name:                    "should fail, if could not get children of Mega repository root",
			loginError:              nil,
			getChildrenError:        errGetChildren,
			getRootNodeHashFunction: nil,
			expRootNodeHash:         "",
			expErr:                  errGetChildren,
		},
		{
			name:                    "should fail, if could not find the expected project root node",
			loginError:              nil,
			getChildrenError:        nil,
			getRootNodeHashFunction: nil,
			expRootNodeHash:         "",
			expErr:                  errGetRootNodeHash,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := mockClient{
				errLogin: test.loginError,
			}
			mockFs := mockFs{
				errGetChildren: test.getChildrenError,
			}
			storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
			if test.getRootNodeHashFunction != nil {
				storageBrowser.getRootNodeHash = test.getRootNodeHashFunction
			}

			err := storageBrowser.Initialize()

			assert.Equal(t, test.expRootNodeHash, storageBrowser.rootNodeHash)
			assert.Equal(t, test.expErr, err)
		})
	}
}

func TestGetRootNodeHash(t *testing.T) {
	tests := []struct {
		name            string
		nodes           []Node
		expRootNodeHash string
		expErr          error
	}{
		{
			name: "should successfully get root node hash, if given list containing a root node",
			nodes: []Node{
				&mockNode{
					name:     rootNodeName,
					nodeType: 1,
					hash:     expRootNodeHash,
				},
			},
			expRootNodeHash: expRootNodeHash,
			expErr:          nil,
		},
		{
			name: "should fail, if none of the nodes is project root",
			nodes: []Node{
				&mockNode{}, &mockNode{},
			},
			expRootNodeHash: "",
			expErr:          errGetRootNodeHash,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rootNodeHash, err := getRootNodeHash(test.nodes)

			assert.Equal(t, test.expRootNodeHash, rootNodeHash)
			assert.Equal(t, test.expErr, err)
		})
	}
}

func TestGetObjectNodeHash(t *testing.T) {
	tests := []struct {
		name                string
		getChildrenError    error
		getChildrenFunction getChildrenFunc
		givenPath           string
		targetSeparator     interface{}
		expHash             string
		expErr              error
	}{
		{
			name:                "should successfully get object node hash, if all inputs are correct",
			getChildrenError:    nil,
			getChildrenFunction: mockGetChildren,
			givenPath:           filepath.Join(expDirName, expFileName),
			targetSeparator:     nil,
			expHash:             expFileHash,
			expErr:              nil,
		},
		{
			name:                "should fail, if could not find expected directory",
			getChildrenError:    nil,
			getChildrenFunction: mockGetChildren,
			givenPath:           filepath.Join("unexpectedDir", expFileName),
			targetSeparator:     nil,
			expHash:             "",
			expErr:              fmt.Errorf("could not find directory: unexpectedDir"),
		},
		{
			name:                "should fail, if could not find expected file",
			getChildrenError:    nil,
			getChildrenFunction: mockGetChildren,
			givenPath:           filepath.Join(expDirName, "unexpFile"),
			targetSeparator:     nil,
			expHash:             "",
			expErr:              fmt.Errorf("could not find file: unexpFile"),
		},
		{
			name:                "should fail, if both path and separator is empty",
			getChildrenError:    nil,
			getChildrenFunction: mockGetChildren,
			givenPath:           "",
			targetSeparator:     "",
			expHash:             "",
			expErr:              fmt.Errorf("could not find object node for "),
		},
		{
			name:                "should fail, if given path is empty",
			getChildrenError:    nil,
			getChildrenFunction: mockGetChildren,
			givenPath:           "",
			targetSeparator:     nil,
			expHash:             "",
			expErr:              fmt.Errorf("trying to find object node for an empty path"),
		},
		{
			name:                "should fail, if could not get children of a node",
			getChildrenError:    errGetChildren,
			getChildrenFunction: nil,
			givenPath:           "path",
			targetSeparator:     nil,
			expHash:             "",
			expErr:              errGetChildren,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockFs := mockFs{
				errGetChildren: test.getChildrenError,
			}
			mockClient := mockClient{}
			storageBrowser := NewMegaBrowser(&mockClient, &mockFs)
			storageBrowser.getRootNodeHash = getRootNodeHash
			if test.getChildrenFunction != nil {
				storageBrowser.getChildren = test.getChildrenFunction
			}
			if test.targetSeparator != nil {
				require.IsType(t, "", test.targetSeparator)
				storageBrowser.targetSeparator = test.targetSeparator.(string)
			}

			result, err := storageBrowser.GetObjectNode(test.givenPath)

			assert.Equal(t, test.expHash, result)
			assert.Equal(t, test.expErr, err)
		})
	}
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

func (m *mockClient) DownloadFile(src *mega.Node, dstpath string, progress *chan int) error {
	return m.errDownload
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
