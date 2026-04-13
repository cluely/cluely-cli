package cmd

import (
	"fmt"
	"os"

	"github.com/cluely/cli/internal/daemon"
	"github.com/spf13/cobra"
)

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if the watch service is running",
	RunE: func(cmd *cobra.Command, args []string) error {
		running, execCmd, err := daemon.Status()
		if err != nil {
			return err
		}

		if !running {
			fmt.Println("Not running.")
			os.Exit(1)
		}

		fmt.Println("Running.")
		if execCmd != "" {
			fmt.Printf("Exec: %s\n", execCmd)
		}
		fmt.Printf("Logs: %s\n", daemon.LogPath())
		return nil
	},
}

func init() {
	daemonCmd.AddCommand(daemonStatusCmd)
}
