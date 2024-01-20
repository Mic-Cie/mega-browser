package megabrowser

import (
	"fmt"

	"github.com/t3rm1n4l/go-mega"
)

const (
	// fileType is the integer value, which specifies that a node is a file, not a directory.
	fileType = 0
	// directoryType is the integer value, which specifies that a node is a directory, not a file.
	directoryType = 1
)

type Node interface {
	GetName() string
	GetType() int
	GetHash() string
	GetSize() int64
}

type getNodeSizeFunc func(Node) int64

// nodeStructArrToInterfaceArr converts array of mega.Node structures to an array of Node interface instances, to make it more generic and allow testing.
func nodeStructArrToInterfaceArr(nodes []*mega.Node) []Node {
	convertedNodes := make([]Node, len(nodes))
	for i, node := range nodes {
		convertedNodes[i] = node
	}
	return convertedNodes
}

// getNodeHashOfExpectedFile expects that given list of nodes contains a file of specific name. Returns that file's hash, otherwise, if that file is not found, returns an error
func getNodeHashOfExpectedFile(expectedFile string, currDirChildNodes *[]Node) (string, error) {
	hash := getNodeHashOfExpectedItem(expectedFile, fileType, currDirChildNodes)
	if hash == "" {
		return "", fmt.Errorf("could not find file: %s", expectedFile)
	}
	return hash, nil
}

// getNodeHashOfExpectedDirectory expects that given list of nodes contains a directory of specific name. Returns that directory's hash, otherwise, if that file is not found, returns an error
func getNodeHashOfExpectedDirectory(expectedDirectory string, currDirChildNodes *[]Node) (string, error) {
	hash := getNodeHashOfExpectedItem(expectedDirectory, directoryType, currDirChildNodes)
	if hash == "" {
		return "", fmt.Errorf("could not find directory: %s", expectedDirectory)
	}
	return hash, nil
}

func getNodeHashOfExpectedItem(expectedItem string, expectedItemType int, currDirChildNodes *[]Node) string {
	for _, child := range *currDirChildNodes {
		if child.GetName() == expectedItem && child.GetType() == expectedItemType {
			return child.GetHash()
		}
	}
	return ""
}

// getNodeSize is a wrapper function that calls the node's Getsize() function. Used for extra abstraction layer.
func getNodeSize(node Node) int64 {
	return node.GetSize()
}
