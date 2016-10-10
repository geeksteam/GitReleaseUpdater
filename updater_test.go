package updater_test

import (
	"errors"
	"testing"

	updater "github.com/geeksteam/GitReleaseUpdater"
)

func TestUpdate(t *testing.T) {
	const currentVersion = "5.0"

	s := source{}

	needUpdate := func(v string) bool {
		return v != currentVersion
	}

	updateMethod := func(path string) error {
		if path != s.path {
			return errors.New("Bad path for update")
		}
		return nil
	}

	restartCommand := func() error {
		return nil
	}

	u := updater.Updater{
		Source:         s,
		NeedUpdate:     needUpdate,
		UpdateMethod:   updateMethod,
		RestartCommand: restartCommand,
	}

	if err := u.Update(); err != nil {
		t.Fatal(err)
	}

}

type source struct {
	version string
	path    string
}

func (s source) LastVersion() (string, error) {
	return s.version, nil
}

func (s source) Download() (string, error) {
	return s.path, nil
}
