package updater

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const defaultDownloadDir = "/tmp"

type DirectLinks struct {
	VersionURL  string
	DownloadURL string

	DownloadDir string
}

func SourceDirectLinks(versionURL, downloadURL string) DirectLinks {
	return DirectLinks{
		VersionURL:  versionURL,
		DownloadURL: downloadURL,
		DownloadDir: defaultDownloadDir,
	}
}

func (dl *DirectLinks) LastVersion() (string, error) {
	resp, err := http.Get(dl.VersionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var b bytes.Buffer

	io.Copy(&b, resp.Body)

	return strings.TrimSpace(b.String()), nil
}

func (dl *DirectLinks) Download() (string, error) {
	resp, err := http.Get(dl.DownloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokens := strings.Split(dl.DownloadURL, "/")
	fileName := filepath.Join(dl.DownloadDir, tokens[len(tokens)-1])

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	io.Copy(file, resp.Body)

	return fileName, nil
}
