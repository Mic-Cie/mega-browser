package megabrowser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/t3rm1n4l/go-mega"
)

var errGetRootNodeHash = fmt.Errorf("failed to get root node hash")

type StorageBrowser interface {
	GetObjectNode(file string) (string, error)
}

type MegaBrowser struct {
	login           string
	pass            string
	rootNodeName    string
	megaClient      StorageClient
	megaFs          Fs
	downloader      Downloader
	getRootNodeHash getRootNodeHashFunc
	getChildren     getChildrenFunc
	targetSeparator string
	rootNodeHash    string
}

type getRootNodeHashFunc func(nodes []Node, rootNodeName string) (string, error)
type getChildrenFunc func(fs Fs, nodeHash string) ([]Node, error)

/*
NewMegaBrowser creates a browser object for a Mega repository.

Expected input parameters are:

	login, pass - credentials for the Mega repository, containing the updated project.
	rootNodeName - name of the directory, containing the updated project.
	megaClient - client for the Mega repository. Can be created with mega.New() function from t3rm1n4l/go-mega package. Make sure to create the megaClient object before actually calling NewMegaBrowser() function.
	fs - system of Mega nodes. FS parameter of the megaClient above can be used.
	downloader - object responsible for downloading and updating project files. Can be created with NewMegaDownloader(client StorageClient) function from this package.
*/
func NewMegaBrowser(login string, pass string, rootNodeName string, megaClient StorageClient, fs Fs, downloader Downloader) *MegaBrowser {
	browser := &MegaBrowser{
		login:           login,
		pass:            pass,
		rootNodeName:    rootNodeName,
		megaClient:      megaClient,
		megaFs:          fs,
		downloader:      downloader,
		getRootNodeHash: getRootNodeHash,
		getChildren:     getChildren,
		targetSeparator: "/",
	}
	return browser
}

/*
Initialize logs in to the Mega repository and initializes the browser parameters, based on that repository.

Returns an error if:

	failed to login
	an error occured while getting children of a repository root node
	could not find the project root node
*/
func (mb *MegaBrowser) Initialize() error {
	err := mb.megaClient.Login(mb.login, mb.pass)
	if err != nil {
		return err
	}

	nodes, err := mb.megaFs.GetChildren(mb.megaFs.GetRoot())
	if err != nil {
		return err
	}

	convertedNodes := nodeStructArrToInterfaceArr(nodes)
	rootNodeHash, err := mb.getRootNodeHash(convertedNodes, mb.rootNodeName)
	if err != nil {
		return err
	}
	mb.rootNodeHash = rootNodeHash

	return nil
}

/*
GetObjectNode takes path to a file and returns its hash from the Mega repository.

Returns an error if:

	given path is empty
	an error occured while getting children of a node
	expected to find node of a file or directory, but did not find it
*/
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

// UpdateFile updates a file at specified localDownloadPath with a file downloaded from Mega node
func (mb *MegaBrowser) UpdateFile(node *mega.Node, localDownloadPath string) error {
	return mb.downloader.DownloadFile(node, localDownloadPath)
}

// getRoodNodeHash takes an array of nodes and checks if any of them is a project root node.
//
// A node is considered a root node, if its name is the same as rootNodeName, and it is a directory.
//
// Returns hash of the node. If could not find the root node, returns an error.
func getRootNodeHash(nodes []Node, rootNodeName string) (string, error) {
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
