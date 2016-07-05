package ctl

import (
	"fmt"
	"os/exec"
)

func update() error {
	if err := pullImages("bryanl/do-autoscale"); err != nil {
		return fmt.Errorf("fetching autoscale image: %s", err)
	}

	if err := exec.Command("systemctl", "restart", "autoscale").Run(); err != nil {
		return fmt.Errorf("updating autoscale servce: %s", err)
	}

	return nil
}
