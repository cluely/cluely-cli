package cmd

import "github.com/spf13/cobra"

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage the background watch service",
	Long:  "Install, start, stop, and monitor the Cluely session watcher as a background service.",
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
