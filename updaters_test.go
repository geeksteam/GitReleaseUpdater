package updater_test

import (
	"os"
	"testing"

	updater "github.com/geeksteam/GitReleaseUpdater"
)

func TestUntarGz(t *testing.T) {
	const (
		archName    = "updaterTest18823"
		checkText   = "tar.gz test"
		unzipedFile = "/tmp/updaterTest18823/test.txt"
	)

	if err := updater.UMUntarGz("/tmp")("./tests/" + archName + ".tar.gz"); err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll("/tmp/" + archName)

	checkFileText(t, unzipedFile, checkText)
}
