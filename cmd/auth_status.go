package cmd

import (
	"fmt"
	"os"

	"github.com/cluely/cli/internal/auth"
	"github.com/spf13/cobra"
)

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  "Display current login status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			return fmt.Errorf("failed to read credentials: %w", err)
		}

		if token == "" {
			fmt.Println("Not logged in. Run 'cluely auth login' to authenticate.")
			os.Exit(1)
		}

		fmt.Println("Logged in.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
}
