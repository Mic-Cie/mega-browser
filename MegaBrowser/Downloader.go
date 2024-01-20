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
	client     StorageClient
	removeFile removeFileFunc
	mkDir      mkdirFunc
	getWd      getWdFunc
}

type removeFileFunc func(path string) error
type mkdirFunc func(path string, perm fs.FileMode) error
type getWdFunc func() (string, error)

func NewMegaDownloader(client StorageClient) *MegaDownloader {
	return &MegaDownloader{
		client:     client,
		removeFile: os.Remove,
		mkDir:      os.MkdirAll,
		getWd:      os.Getwd,
	}
}

func (md *MegaDownloader) DownloadFile(node *mega.Node, localDownloadPath string) error {
	err := md.removeOutdatedFile(localDownloadPath)
	if err != nil {
		return err
	}

	err = md.createFileDirectoryIfNotExist(localDownloadPath)
	if err != nil {
		return err
	}

	currentDir, err := md.getWd()
	if err != nil {
		return err
	}

	err = md.client.DownloadFile(node, filepath.Join(currentDir, localDownloadPath), nil)
	if err != nil {
		return err
	}

	return nil
}

// removeOutdatedFile removes a file that is supposed to be updated.
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

// createFileDirectoryIfNotExist creates path for a file that is supposed to be downloaded. Does nothing, if the directory already exists.
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
