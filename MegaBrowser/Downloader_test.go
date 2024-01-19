package megabrowser

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errRemoveFile = fmt.Errorf("mock remove file error")

func TestShouldThrowNoErrorWhenDownloadingFileThatNotExistLocally(t *testing.T) {
	downloader := NewMegaDownloader()

	err := downloader.DownloadFile(nil, "path/that/not/exist")

	assert.Nil(t, err)
}

func TestShouldThrowNoErrorWhenDownloadingFileThatExistedLocally(t *testing.T) {
	downloader := NewMegaDownloader()
	downloader.removeFile = mockRemoveFileSuccess

	err := downloader.DownloadFile(nil, filepath.Join("testDir", "localFile.txt"))

	assert.Nil(t, err)
}

func TestShouldFailIfGivenPathisIncorrect(t *testing.T) {
	downloader := NewMegaDownloader()
	path := strings.Repeat("?", 1000)

	err := downloader.DownloadFile(nil, path)

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), path)
}

func TestShouldThrowNoErrorWhenFailedToRemoveALocalFile(t *testing.T) {
	downloader := NewMegaDownloader()
	downloader.removeFile = mockRemoveFileFail

	err := downloader.DownloadFile(nil, filepath.Join("testDir", "localFile.txt"))

	assert.Equal(t, errRemoveFile, err)
}

func mockRemoveFileSuccess(path string) error {
	return nil
}

func mockRemoveFileFail(path string) error {
	return errRemoveFile
}
