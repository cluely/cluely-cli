package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/api"
	"github.com/spf13/cobra"
)

var sessionsTagCmd = &cobra.Command{
	Use:     "tag <session-id> <tag-id>",
	Short:   "Add a tag to a session",
	Args:    cobra.ExactArgs(2),
	Example: `  cluely sessions tag <session-id> <tag-id>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := api.Call("sessionTags/create", map[string]string{
			"sessionId": args[0],
			"tagId":     args[1],
		}, nil); err != nil {
			return err
		}

		fmt.Println("Tag added.")
		return nil
	},
}

var sessionsUntagCmd = &cobra.Command{
	Use:     "untag <session-id> <tag-id>",
	Short:   "Remove a tag from a session",
	Args:    cobra.ExactArgs(2),
	Example: `  cluely sessions untag <session-id> <tag-id>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := api.Call("sessionTags/delete", map[string]string{
			"sessionId": args[0],
			"tagId":     args[1],
		}, nil); err != nil {
			return err
		}

		fmt.Println("Tag removed.")
		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsTagCmd)
	sessionsCmd.AddCommand(sessionsUntagCmd)
}
