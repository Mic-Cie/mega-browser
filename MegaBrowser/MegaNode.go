package megabrowser

import "github.com/t3rm1n4l/go-mega"

type Node interface {
	GetName() string
	GetType() int
	GetHash() string
}

func nodeStructArrToInterfaceArr(nodes []*mega.Node) []Node {
	convertedNodes := make([]Node, len(nodes))
	for i, node := range nodes {
		convertedNodes[i] = node
	}
	return convertedNodes
}
