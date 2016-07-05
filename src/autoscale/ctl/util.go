package ctl

import (
	"fmt"
	"os/exec"
)

func pullImages(images ...string) error {
	for _, image := range images {
		fmt.Println("* pulling", image, "image")
		if err := exec.Command("/usr/bin/docker", "pull", image).Run(); err != nil {
			return err
		}
	}

	return nil
}
