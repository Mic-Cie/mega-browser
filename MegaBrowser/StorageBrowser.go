package megabrowser

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	// login of the Mega repository owner.
	login = `login`
	// password of the Mega repository owner.
	pass = `password`
	// rootNodeName is the name of the directory which will be considered root directory of the updated application. Make sure to place it in the root of your Mega repository.
	rootNodeName = "root-node-name"
)

var errGetRootNodeHash = fmt.Errorf("failed to get root node hash")

type StorageBrowser interface {
	GetObjectNode(file string) (string, error)
}

type MegaBrowser struct {
	megaClient      StorageClient
	megaFs          Fs
	getRootNodeHash getRootNodeHashFunc
	getChildren     getChildrenFunc
	targetSeparator string
	rootNodeHash    string
}

type getRootNodeHashFunc func(nodes []Node) (string, error)
type getChildrenFunc func(fs Fs, nodeHash string) ([]Node, error)

func NewMegaBrowser(megaClient StorageClient, fs Fs) *MegaBrowser {
	browser := &MegaBrowser{
		megaClient:      megaClient,
		megaFs:          fs,
		getRootNodeHash: getRootNodeHash,
		getChildren:     getChildren,
		targetSeparator: "/",
	}
	return browser
}

// Initialize logs in to the Mega repository and initializes the browser parameters, based on that repository.
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

func (mb *MegaBrowser) GetObjectNode(file string) (string, error) {
	splitPath := strings.Split(filepath.ToSlash(file), mb.targetSeparator)
	len := len(splitPath)
	if len == 1 && splitPath[0] == "" {
		return "", fmt.Errorf("trying to find object node for an empty path")
	}

	var targetFile string
	var currentDir string
	if len != 0 {
		targetFile = splitPath[len-1]
		currentDir = mb.rootNodeHash
	}

	for i, _ := range splitPath {
		childNodes, err := mb.getChildren(mb.megaFs, currentDir)
		if err != nil {
			return "", err
		}

		if i == len-1 {
			result, err := getNodeHashOfExpectedFile(targetFile, &childNodes)
			if err != nil {
				return "", err
			}
			return result, nil
		} else {
			var err error
			currentDir, err = getNodeHashOfExpectedDirectory(splitPath[i], &childNodes)
			if err != nil {
				return "", err
			}
		}
	}
	return "", fmt.Errorf("could not find object node for %s", file)
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

func getChildren(fs Fs, nodeHash string) ([]Node, error) {
	currentDirNode := fs.HashLookup(nodeHash)
	currDirChildNodes, err := fs.GetChildren(currentDirNode)
	if err != nil {
		return nil, err
	}
	return nodeStructArrToInterfaceArr(currDirChildNodes), nil
}
