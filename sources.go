package updater

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	return downloadFrom(dl.DownloadURL, dl.DownloadDir)
}

const githubAPIURL = "https://api.github.com"

type GitReleases struct {
	// Credentials is optional and uses only for private repos
	Username string
	Password string

	// Repository info
	Owner string
	Repo  string

	DownloadURL string
	DownloadDir string
}

func (gr *GitReleases) LastVersion() (string, error) {
	type release struct {
		Version     string `json:"name"`
		DownloadURL string `json:"tarball_url"`
		Description string `json:"body"`
	}

	lastReleaseURL := fmt.Sprintf("%s/repos/%s/%s/releases/latest", githubAPIURL, gr.Owner, gr.Repo)

	req, err := http.NewRequest(http.MethodGet, lastReleaseURL, nil)
	if err != nil {
		return "", err
	}

	if gr.Username != "" && gr.Password != "" {
		req.SetBasicAuth(gr.Username, gr.Password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var rel release

	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}

	if gr.DownloadURL == "" {
		gr.DownloadURL = rel.DownloadURL
	}

	return rel.Version, nil
}

func (gr *GitReleases) Download() (string, error) {
	return downloadFrom(gr.DownloadURL, gr.DownloadDir)
}

func downloadFrom(url, toDir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokens := strings.Split(url, "/")
	fileName := filepath.Join(toDir, tokens[len(tokens)-1])

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil && err != io.EOF {
		return "", err
	}

	return fileName, nil
}
