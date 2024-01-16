package megabrowser

import (
	"fmt"
)

const (
	// login of the Mega repository owner.
	login = `login`
	// password of the Mega repository owner.
	pass = `password`
	// rootNodeName is the name of the directory which will be considered root directory of the updated application. Make sure to place it in the root of your Mega repository.
	rootNodeName = "root-node-name"
	// directoryType is the integer value, which specifies that node is a directory, not a file.
	directoryType = 1
)

var errGetRootNodeHash = fmt.Errorf("failed to get root node hash")

type StorageBrowser interface {
	GetObjectNode(file string) (string, error)
}

type MegaBrowser struct {
	megaClient      StorageClient
	megaFs          Fs
	getRootNodeHash getRootNodeHashFunc
	rootNodeHash    string
}

type getRootNodeHashFunc func(nodes []Node) (string, error)

func NewMegaBrowser(megaClient StorageClient, fs Fs) *MegaBrowser {
	browser := &MegaBrowser{
		megaClient:      megaClient,
		megaFs:          fs,
		getRootNodeHash: getRootNodeHash,
	}
	return browser
}

func (mb *MegaBrowser) Initialize() error {
	err := mb.megaClient.Login(login, pass)
	if err != nil {
		return err
	}

	nodes, err := mb.megaFs.GetChildren(mb.megaFs.GetRoot())
	if err != nil {
		return err
	}

	convertedNodes := nodeStructArrToInterfaceArr(nodes)
	rootNodeHash, err := mb.getRootNodeHash(convertedNodes)
	if err != nil {
		return err
	}
	mb.rootNodeHash = rootNodeHash

	return nil
}

func (b *MegaBrowser) GetObjectNode(file string) (string, error) {
	return "", nil
}

func getRootNodeHash(nodes []Node) (string, error) {
	for _, node := range nodes {
		nodeName := node.GetName()
		nodeType := node.GetType()
		if nodeName == rootNodeName && nodeType == directoryType {
			rootNodeHash := node.GetHash()
			return rootNodeHash, nil
		}
	}
	return "", errGetRootNodeHash
}
