package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/auth"
	"github.com/cluely/cli/internal/config"
	"github.com/spf13/cobra"
)

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Cluely",
	Long:  "Authenticate with Cluely by opening a browser window.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if auth.HasToken() {
			fmt.Println("Already logged in. Use 'cluely auth logout' first to switch accounts.")
			return nil
		}

		if err := auth.Login(cmd.Context(), config.WebURL); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Println("Successfully logged in!")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
}
