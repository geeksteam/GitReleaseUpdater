package updater

import (
	"fmt"
	"os/exec"
)

func RCService(service string) func() error {
	return func() error {
		if err := exec.Command("service", service, "restart").Run(); err != nil {
			return fmt.Errorf("Failed to restart service %s: %s", service, err.(*exec.ExitError).Stderr)
		}

		return nil
	}
}
