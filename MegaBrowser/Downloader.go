package megabrowser

import (
	"os"

	"github.com/t3rm1n4l/go-mega"
)

type Downloader interface {
	DownloadFile(node *mega.Node, localDownloadPath string) error
}

type MegaDownloader struct {
	removeFile removeFileFunc
}

type removeFileFunc func(path string) error

func NewMegaDownloader() *MegaDownloader {
	return &MegaDownloader{
		removeFile: os.Remove,
	}
}

func (md *MegaDownloader) DownloadFile(node *mega.Node, localDownloadPath string) error {
	err := md.removeOutdatedFile(localDownloadPath)
	if err != nil {
		return err
	}

	return nil
}

func (md *MegaDownloader) removeOutdatedFile(localDownloadPath string) error {
	if _, err := os.Stat(localDownloadPath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err = md.removeFile(localDownloadPath)
		if err != nil {
			return err
		}
	}
	return nil
}
