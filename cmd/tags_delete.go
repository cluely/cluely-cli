package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/api"
	"github.com/spf13/cobra"
)

var tagsDeleteCmd = &cobra.Command{
	Use:   "delete <tag-id>",
	Short: "Delete a tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tagID := args[0]

		if err := api.Call("tags/delete", map[string]string{"id": tagID}, nil); err != nil {
			return err
		}

		fmt.Println("Tag deleted.")
		return nil
	},
}

func init() {
	tagsCmd.AddCommand(tagsDeleteCmd)
}
