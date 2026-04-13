package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cluely/cli/internal/api"
	"github.com/cluely/cli/internal/color"
	"github.com/spf13/cobra"
)

var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonOut, _ := cmd.Flags().GetBool("json")

		if jsonOut {
			raw, err := api.CallRaw("tags/list", nil)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		}

		var tags []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		}

		if err := api.Call("tags/list", nil, &tags); err != nil {
			return err
		}

		if len(tags) == 0 {
			fmt.Println("No tags found. Create one with 'cluely tags create'.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tNAME\tPREVIEW\n")
		for _, t := range tags {
			fmt.Fprintf(w, "%s\t%s\t%s\n", t.ID, t.Name, color.TagBadge(t.Name, t.Color))
		}
		w.Flush()

		return nil
	},
}

func init() {
	tagsCmd.AddCommand(tagsListCmd)
}
