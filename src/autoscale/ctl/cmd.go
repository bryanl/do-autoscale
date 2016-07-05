package ctl

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// RootCmd is the root command for autoscalectl.
	RootCmd = &cobra.Command{
		Use:   "autoscalectl",
		Short: "autoscalectl controls the autoscaler",
	}
)

// InitCommands initializes the commands.
func InitCommands() {
	setupCmd := &cobra.Command{
		Use: "setup",
		Run: func(cmd *cobra.Command, args []string) {
			if err := setup(); err != nil {
				logrus.WithError(err).Fatal("unabe to run setup")
			}
		},
	}

	RootCmd.AddCommand(setupCmd)
}
