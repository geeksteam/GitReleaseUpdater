package updater

import (
	"errors"
	"time"

	"github.com/geeksteam/klogger"
)

type Source interface {
	LastVersion() (string, error)
	Download() (string, error)
}

type Updater struct {
	Logger klogger.Logger

	NeedUpdate     func(string) bool
	UpdateMethod   func(string) error
	RestartCommand func() error

	Source Source

	CheckInterval time.Duration
}

func (u Updater) Update() error {
	if err := u.checkDependencies(); err != nil {
		return err
	}

	lastVer, err := u.Source.LastVersion()
	if err != nil {
		return err
	}

	if !u.NeedUpdate(lastVer) {
		return nil
	}

	if u.Logger != nil {
		u.Logger.Info("Found new version %s, downloading", lastVer)
	}

	latestPath, err := u.Source.Download()
	if err != nil {
		return err
	}

	if err := u.UpdateMethod(latestPath); err != nil {
		return err
	}

	if err := u.RestartCommand(); err != nil {
		return err
	}

	return nil
}

func (u Updater) checkDependencies() error {
	if u.NeedUpdate == nil {
		return errors.New("Missed NeedUpdate function")
	}

	if u.UpdateMethod == nil {
		return errors.New("Missed UpdateMethod function")
	}

	if u.RestartCommand == nil {
		return errors.New("Missed RestartCommand function")
	}

	if u.Source == nil {
		return errors.New("Missed update Source")
	}

	return nil
}
