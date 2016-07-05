package main

import (
	"autoscale/ctl"
	"os"

	"github.com/Sirupsen/logrus"
)

func main() {
	ctl.InitCommands()
	if err := ctl.RootCmd.Execute(); err != nil {
		logrus.WithError(err).Error("command failed")
		os.Exit(-1)
	}
}
