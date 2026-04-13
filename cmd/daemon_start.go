package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/daemon"
	"github.com/spf13/cobra"
)

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the background watch service",
	Long: `Install and start the Cluely session watcher as a background service.

The service runs on login and auto-restarts on failure.
Uses launchd on macOS and systemd on Linux.`,
	Example: `  cluely daemon start --exec "echo \$CLUELY_SESSION_TITLE finished"
  cluely daemon start --exec "./on-session-complete.sh"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		execCmd, _ := cmd.Flags().GetString("exec")
		if execCmd == "" {
			return fmt.Errorf("--exec is required")
		}

		if err := daemon.Install(execCmd); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}

		fmt.Println("Service started.")
		fmt.Printf("Logs: %s\n", daemon.LogPath())
		return nil
	},
}

func init() {
	daemonStartCmd.Flags().String("exec", "", "Command to run when a session finishes (required)")
	daemonStartCmd.MarkFlagRequired("exec")
	daemonCmd.AddCommand(daemonStartCmd)
}
