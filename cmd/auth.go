package cmd

import "github.com/spf13/cobra"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "Log in, log out, and check authentication status.",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
