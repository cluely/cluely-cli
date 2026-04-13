package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cluely/cli/internal/daemon"
	"github.com/spf13/cobra"
)

var daemonLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Tail the watch service logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		logPath := daemon.LogPath()

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			return fmt.Errorf("no log file found at %s", logPath)
		}

		follow, _ := cmd.Flags().GetBool("follow")

		var tail *exec.Cmd
		if follow {
			tail = exec.Command("tail", "-f", logPath)
		} else {
			tail = exec.Command("tail", "-100", logPath)
		}
		tail.Stdout = os.Stdout
		tail.Stderr = os.Stderr
		return tail.Run()
	},
}

func init() {
	daemonLogsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	daemonCmd.AddCommand(daemonLogsCmd)
}
