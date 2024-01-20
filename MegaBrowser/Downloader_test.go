package megabrowser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errRemoveFile = fmt.Errorf("mock remove file error")
	errMkDir      = fmt.Errorf("mock mkdir error")
)

func TestShouldThrowNoErrorWhenDownloadingFileThatNotExistLocally(t *testing.T) {
	downloader := NewMegaDownloader()

	err := downloader.DownloadFile(nil, "temp/path/that/not/exist.txt")

	assert.Nil(t, err)

	dir, err := os.Getwd()
	require.Nil(t, err)
	err = os.Remove(filepath.Join(dir, "temp"))
	require.Nil(t, err)
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

func TestShouldFailWhenFailedToRemoveALocalFile(t *testing.T) {
	downloader := NewMegaDownloader()
	downloader.removeFile = mockRemoveFileFail

	err := downloader.DownloadFile(nil, filepath.Join("testDir", "localFile.txt"))

	assert.Equal(t, errRemoveFile, err)
}

func TestShouldCreateDirectoryIfItNotExists(t *testing.T) {
	downloader := NewMegaDownloader()

	err := downloader.createFileDirectoryIfNotExist("temp/file.txt")

	assert.Nil(t, err)

	dir, err := os.Getwd()
	require.Nil(t, err)
	err = os.Remove(filepath.Join(dir, "temp"))
	require.Nil(t, err)
}

func TestShouldNotCreateDirectoryIfItExists(t *testing.T) {
	downloader := NewMegaDownloader()

	err := downloader.createFileDirectoryIfNotExist("testDir/localFile.txt")

	assert.Nil(t, err)
}

func TestShouldFailIfCouldNotCreateDirectory(t *testing.T) {
	downloader := NewMegaDownloader()
	downloader.mkDir = mockMkDirFail

	err := downloader.createFileDirectoryIfNotExist("temp/file.txt")

	assert.Equal(t, errMkDir, err)
}

func TestShouldFailCreatingDirectoryIfPathIsTooLong(t *testing.T) {
	downloader := NewMegaDownloader()
	path := filepath.Join(strings.Repeat("?", 1000), "path")

	err := downloader.createFileDirectoryIfNotExist(path)

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "???????????????")
}

func mockRemoveFileSuccess(path string) error {
	return nil
}

func mockRemoveFileFail(path string) error {
	return errRemoveFile
}

func mockMkDirFail(path string, perm fs.FileMode) error {
	return errMkDir
}
