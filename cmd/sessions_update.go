package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/api"
	"github.com/spf13/cobra"
)

var sessionsUpdateCmd = &cobra.Command{
	Use:   "update <session-id>",
	Short: "Update a session",
	Long:  "Update the title or summary of a session.",
	Args:  cobra.ExactArgs(1),
	Example: `  cluely sessions update <session-id> --title "Q2 Planning"
  cluely sessions update <session-id> --summary "Discussed roadmap priorities"
  cluely sessions update <session-id> --title "Standup" --summary "Quick sync"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		title, _ := cmd.Flags().GetString("title")
		summary, _ := cmd.Flags().GetString("summary")

		if title == "" && summary == "" {
			return fmt.Errorf("provide at least one of --title or --summary")
		}

		input := map[string]interface{}{"id": sessionID}
		if title != "" {
			input["title"] = title
		}
		if summary != "" {
			input["summary"] = summary
		}

		if err := api.Call("sessions/update", input, nil); err != nil {
			return err
		}

		fmt.Println("Session updated.")
		return nil
	},
}

func init() {
	sessionsUpdateCmd.Flags().String("title", "", "New session title")
	sessionsUpdateCmd.Flags().String("summary", "", "New session summary")
	sessionsCmd.AddCommand(sessionsUpdateCmd)
}
