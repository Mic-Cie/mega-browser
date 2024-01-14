package megabrowser

import "github.com/t3rm1n4l/go-mega"

type StorageClient interface {
	Login(login string, pass string) error
	GetFS() *mega.MegaFS
}

type MegaClient struct {
	mega.Mega
}

func (m *MegaClient) GetFS() *mega.MegaFS {
	return m.FS
}
