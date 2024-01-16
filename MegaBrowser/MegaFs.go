package megabrowser

import "github.com/t3rm1n4l/go-mega"

type Fs interface {
	GetChildren(*mega.Node) ([]*mega.Node, error)
	GetRoot() *mega.Node
}
