package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/daemon"
	"github.com/spf13/cobra"
)

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the background watch service",
	Long:  "Stop and uninstall the Cluely session watcher service.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.Uninstall(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}

		fmt.Println("Service stopped and removed.")
		return nil
	},
}

func init() {
	daemonCmd.AddCommand(daemonStopCmd)
}
