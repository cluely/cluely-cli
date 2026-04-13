package cmd

import (
	"fmt"

	"github.com/cluely/cli/internal/api"
	"github.com/cluely/cli/internal/color"
	"github.com/spf13/cobra"
)

var tagsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a tag",
	Args:  cobra.ExactArgs(1),
	Example: `  cluely tags create "Sales Call" --color "#4f46e5"
  cluely tags create "Interview" --color "#059669"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		tagColor, _ := cmd.Flags().GetString("color")

		var result struct {
			ID string `json:"id"`
		}

		if err := api.Call("tags/create", map[string]string{
			"name":  name,
			"color": tagColor,
		}, &result); err != nil {
			return err
		}

		fmt.Printf("Created tag: %s %s\n", color.TagBadge(name, tagColor), result.ID)
		return nil
	},
}

func init() {
	tagsCreateCmd.Flags().String("color", "#6b7280", "Tag color as hex (e.g. #4f46e5)")
	tagsCmd.AddCommand(tagsCreateCmd)
}
