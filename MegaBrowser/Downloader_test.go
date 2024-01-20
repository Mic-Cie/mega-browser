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
	errGetwd      = fmt.Errorf("mock getwd error")
)

func TestDownloadFileSuccessCase(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		removeFileFunction removeFileFunc
		needCleanup        bool
	}{
		{
			name:               "should not fail, if downloading file that does not exist locally",
			path:               "temp/path/that/not/exist.txt",
			removeFileFunction: nil,
			needCleanup:        true,
		},
		{
			name:               "should not fail, if downloading file that does exist locally",
			path:               filepath.Join("testDir", "localFile.txt"),
			removeFileFunction: mockRemoveFileSuccess,
			needCleanup:        false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			downloader := NewMegaDownloader(
				&mockClient{},
			)
			if test.removeFileFunction != nil {
				downloader.removeFile = test.removeFileFunction
			}

			err := downloader.DownloadFile(nil, test.path)
			if test.needCleanup {
				defer cleanupTestDir(t)
			}

			assert.Nil(t, err)
		})
	}
}

func TestDownloadFileFailCase(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		removeFileFunction removeFileFunc
		mkdirFunction      mkdirFunc
		getWdFunction      getWdFunc
		downloadErr        error
		expErr             string
	}{
		{
			name:               "should fail, if given path is incorrect",
			path:               strings.Repeat("?", 1000),
			removeFileFunction: nil,
			mkdirFunction:      nil,
			downloadErr:        nil,
			getWdFunction:      nil,
			expErr:             strings.Repeat("?", 1000),
		},
		{
			name:               "should fail, if could not remove the local file",
			path:               filepath.Join("testDir", "localFile.txt"),
			removeFileFunction: mockRemoveFileFail,
			mkdirFunction:      nil,
			downloadErr:        nil,
			getWdFunction:      nil,
			expErr:             errRemoveFile.Error(),
		},
		{
			name:               "should fail, if could not create a directory",
			path:               "temp/path/that/not/exist.txt",
			removeFileFunction: nil,
			mkdirFunction:      mockMkDirFail,
			downloadErr:        nil,
			getWdFunction:      nil,
			expErr:             errMkDir.Error(),
		},
		{
			name:               "should fail, if could not download file",
			path:               filepath.Join("testDir", "localFile.txt"),
			removeFileFunction: mockRemoveFileSuccess,
			mkdirFunction:      nil,
			downloadErr:        errDownload,
			getWdFunction:      nil,
			expErr:             errDownload.Error(),
		},
		{
			name:               "should fail, if could not get working directory",
			path:               filepath.Join("testDir", "localFile.txt"),
			removeFileFunction: mockRemoveFileSuccess,
			mkdirFunction:      nil,
			downloadErr:        nil,
			getWdFunction:      mockGetWdFail,
			expErr:             errGetwd.Error(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &mockClient{
				errDownload: test.downloadErr,
			}
			downloader := NewMegaDownloader(
				client,
			)
			if test.removeFileFunction != nil {
				downloader.removeFile = test.removeFileFunction
			}
			if test.mkdirFunction != nil {
				downloader.mkDir = test.mkdirFunction
			}
			if test.getWdFunction != nil {
				downloader.getWd = test.getWdFunction
			}

			err := downloader.DownloadFile(nil, test.path)

			require.NotNil(t, err)
			require.NotEmpty(t, test.expErr)
			assert.Contains(t, err.Error(), test.expErr)
		})
	}
}

func TestCreateDirectoryIfItNotExistsSuccessCase(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		needCleanup bool
	}{
		{
			name:        "should not fail, if target directory does not exist",
			path:        "temp/file.txt",
			needCleanup: true,
		},
		{
			name:        "should not fail, if target directory already exists",
			path:        "testDir/localFile.txt",
			needCleanup: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			downloader := NewMegaDownloader(
				&mockClient{},
			)
			if test.needCleanup {
				defer cleanupTestDir(t)
			}

			err := downloader.createFileDirectoryIfNotExist(test.path)

			assert.Nil(t, err)
		})
	}
}

func TestCreateDirectoryIfItNotExistsFailCase(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		mkdirFunction mkdirFunc
		expErr        string
	}{
		{
			name:          "should fail, if could not create the directory",
			path:          "temp/file.txt",
			mkdirFunction: mockMkDirFail,
			expErr:        errMkDir.Error(),
		},
		{
			name:          "should fail, if given path is incorrect",
			path:          filepath.Join(strings.Repeat("?", 1000), "path"),
			mkdirFunction: nil,
			expErr:        strings.Repeat("?", 1000),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			downloader := NewMegaDownloader(
				&mockClient{},
			)
			if test.mkdirFunction != nil {
				downloader.mkDir = test.mkdirFunction
			}

			err := downloader.createFileDirectoryIfNotExist(test.path)

			require.NotNil(t, err)
			require.NotEmpty(t, test.expErr)
			assert.Contains(t, err.Error(), test.expErr)
		})
	}
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

func mockGetWdFail() (string, error) {
	return "", errGetwd
}

func cleanupTestDir(t *testing.T) {
	dir, err := os.Getwd()
	require.Nil(t, err)
	err = os.RemoveAll(filepath.Join(dir, "temp"))
	require.Nil(t, err)
}
