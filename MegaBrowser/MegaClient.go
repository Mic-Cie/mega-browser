package megabrowser

import "github.com/t3rm1n4l/go-mega"

type StorageClient interface {
	Login(login string, pass string) error
}

type MegaClient struct {
	mega.Mega
}
