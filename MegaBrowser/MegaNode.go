package megabrowser

import "github.com/t3rm1n4l/go-mega"

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
