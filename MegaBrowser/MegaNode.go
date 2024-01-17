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
}

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
	for _, child := range *currDirChildNodes {
		if child.GetName() == expectedFile && child.GetType() == fileType {
			return child.GetHash(), nil
		}
	}
	return "", fmt.Errorf("could not find file: %s", expectedFile)
}
