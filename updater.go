package updater

import (
	"errors"
	"time"

	"github.com/geeksteam/klogger"
)

// DefaultCheckInterval duration which will be used for Updater.Watch method
// when Updater.CheckInterval is 0
var DefaultCheckInterval = 24 * time.Hour

// Source It's entity which knows where to find latest version and how download it
type Source interface {
	LastVersion() (string, error)
	Download() (string, error)
}

// Updater provide main functionality and contains main update flow
type Updater struct {
	Logger klogger.Logger

	NeedUpdate     func(string) bool
	UpdateMethod   func(string) error
	RestartCommand func() error

	Source Source

	CheckInterval time.Duration
}

// Watch Allows to start watching for updates with period = Updater.CheckInterval
// when Updater.CheckInterval is missed DefaultCheckInterval will be used
func (u Updater) Watch() {
	if u.CheckInterval == 0 {
		u.CheckInterval = DefaultCheckInterval
	}

	u.watch()
}

func (u Updater) watch() {
	if u.Logger != nil {
		u.Logger.Info("Looking for updates")
	}

	if err := u.Update(); err != nil {
		if u.Logger != nil {
			u.Logger.Error("Update failed: %s", err)
		}
	}

	time.AfterFunc(u.CheckInterval, u.watch)
}

// Update Looks for new version and if it different from current then downloads
// and installs it. After all calls restart func
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
