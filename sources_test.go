package updater_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	updater "github.com/geeksteam/GitReleaseUpdater"
)

func TestDirectLinks(t *testing.T) {
	const (
		version     = "1.2b"
		dowloadPath = "./latest.test"
		versionURI  = "/latestver"
		downloadURI = "/latest"

		testText = "Text in downloaded file"
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == versionURI {
			fmt.Fprintln(w, version)
			return
		}
		if r.RequestURI == downloadURI {
			fmt.Fprint(w, testText)
		}
	}))
	defer ts.Close()

	dl := updater.SourceDirectLinks(ts.URL+versionURI, ts.URL+downloadURI)

	// Testring versions
	respVer, err := dl.LastVersion()
	if err != nil {
		t.Fatal(err)
	}
	if respVer != version {
		t.Fatalf("Version is %s, must be %s", respVer, version)
	}

	// Testing downloads
	path, err := dl.Download()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	checkFileText(t, path, testText)
}

func checkFileText(t *testing.T, path, text string) {
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var b bytes.Buffer

	io.Copy(&b, file)

	if strings.TrimSpace(b.String()) != text {
		t.Fatalf("Test text in file %s is '%s' instead of '%s'", path, b.String(), text)
	}
}
