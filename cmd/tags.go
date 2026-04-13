package cmd

import "github.com/spf13/cobra"

var tagsCmd = &cobra.Command{
	Use:     "tags",
	Aliases: []string{"t"},
	Short:   "Manage tags",
	Long:    "Create, list, and delete tags for organizing sessions.",
}

func init() {
	tagsCmd.PersistentFlags().Bool("json", false, "Output raw JSON")
	rootCmd.AddCommand(tagsCmd)
}
