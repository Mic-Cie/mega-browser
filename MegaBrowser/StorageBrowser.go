package megabrowser

import "github.com/t3rm1n4l/go-mega"

type StorageBrowser interface {
	GetObjectNode(file string) (string, error)
}

type MegaBrowser struct {
	MegaClient *mega.Mega
}

func NewMegaBrowser() (*MegaBrowser, error) {
	return &MegaBrowser{}, nil
}

func (b *MegaBrowser) GetObjectNode(file string) (string, error) {
	return "", nil
}
