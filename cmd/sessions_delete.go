package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/api"
	"github.com/spf13/cobra"
)

var sessionsDeleteCmd = &cobra.Command{
	Use:   "delete <session-id>",
	Short: "Delete a session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := api.Call("sessions/delete", map[string]string{"id": args[0]}, nil); err != nil {
			return err
		}

		fmt.Println("Session deleted.")
		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsDeleteCmd)
}
