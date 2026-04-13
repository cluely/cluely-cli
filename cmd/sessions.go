package cmd

import "github.com/spf13/cobra"

var sessionsCmd = &cobra.Command{
	Use:     "sessions",
	Aliases: []string{"s"},
	Short:   "Manage sessions",
	Long:    "List, view, and manage meeting sessions.",
}

func init() {
	sessionsCmd.PersistentFlags().Bool("json", false, "Output raw JSON")
	rootCmd.AddCommand(sessionsCmd)
}
