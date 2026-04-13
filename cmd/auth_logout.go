package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/auth"
	"github.com/spf13/cobra"
)

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of Cluely",
	Long:  "Remove stored authentication credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !auth.HasToken() {
			fmt.Println("Not logged in.")
			return nil
		}

		if err := auth.ClearToken(); err != nil {
			return fmt.Errorf("failed to clear credentials: %w", err)
		}

		fmt.Println("Logged out successfully.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLogoutCmd)
}
