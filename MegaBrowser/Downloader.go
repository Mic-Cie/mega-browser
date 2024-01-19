package megabrowser

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/t3rm1n4l/go-mega"
)

type Downloader interface {
	DownloadFile(node *mega.Node, localDownloadPath string) error
}

type MegaDownloader struct {
	removeFile removeFileFunc
	mkDir      mkdirFunc
}

type removeFileFunc func(path string) error
type mkdirFunc func(path string, perm fs.FileMode) error

func NewMegaDownloader() *MegaDownloader {
	return &MegaDownloader{
		removeFile: os.Remove,
		mkDir:      os.MkdirAll,
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

func (md *MegaDownloader) createFileDirectoryIfNotExist(downloadPathWithFile string) error {
	fileDirectory := extractDirectoryFromFullPath(downloadPathWithFile)
	_, err := os.Stat(fileDirectory)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return nil
	}
	err = md.mkDir(fileDirectory, 0777)
	if err != nil {
		return err
	}
	return nil
}

func extractDirectoryFromFullPath(fullPath string) string {
	return filepath.Dir(fullPath)
}
