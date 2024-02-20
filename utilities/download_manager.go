package utilities

import (
	"fmt"
	"io"
	"net/http"
)

type DownloadManager struct {
}

func NewDownloadManager() *DownloadManager {
	return &DownloadManager{}
}

func (dm *DownloadManager) DownloadFile(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file, status code: %d", response.StatusCode)
	}

	fileBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}
