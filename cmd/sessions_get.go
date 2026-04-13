package cmd

import (
	"fmt"
	"strings"

	"github.com/cluely/cli/internal/api"
	"github.com/cluely/cli/internal/color"
	"github.com/spf13/cobra"
)

var allGetSections = []string{"id", "title", "state", "created", "ended", "tags", "attendees", "summary", "transcript"}

var sessionsGetCmd = &cobra.Command{
	Use:   "get <session-id>",
	Short: "Get session details",
	Long:  "Display details and transcript for a specific session.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		jsonOut, _ := cmd.Flags().GetBool("json")
		fields, _ := cmd.Flags().GetString("fields")
		noFields, _ := cmd.Flags().GetString("no-fields")

		input := map[string]string{"id": sessionID}

		if jsonOut {
			raw, err := api.CallRaw("sessions/get", input)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		}

		sections := resolveGetSections(fields, noFields)
		show := func(s string) bool { return sections[s] }

		var result struct {
			ID        string       `json:"id"`
			State     string       `json:"state"`
			Title     *string      `json:"title"`
			Summary   *string      `json:"summary"`
			Tags      []sessionTag `json:"tags"`
			CreatedAt string       `json:"createdAt"`
			EndedAt   *string      `json:"endedAt"`
			Attendees []struct {
				Email string `json:"email"`
			} `json:"attendees"`
			Transcript []struct {
				CreatedAt  string `json:"createdAt"`
				RelativeMs int    `json:"relativeMs"`
				Role       string `json:"role"`
				Text       string `json:"text"`
			} `json:"transcript"`
		}

		if err := api.Call("sessions/get", input, &result); err != nil {
			return err
		}

		if show("id") {
			fmt.Printf("Session: %s\n", result.ID)
		}
		if show("title") {
			title := "(untitled)"
			if result.Title != nil {
				title = *result.Title
			}
			fmt.Printf("Title:   %s\n", title)
		}
		if show("state") {
			fmt.Printf("State:   %s\n", result.State)
		}
		if show("created") {
			fmt.Printf("Created: %s\n", formatTime(result.CreatedAt))
		}
		if show("ended") && result.EndedAt != nil {
			fmt.Printf("Ended:   %s\n", formatTime(*result.EndedAt))
		}
		if show("tags") && len(result.Tags) > 0 {
			badges := make([]string, len(result.Tags))
			for i, t := range result.Tags {
				badges[i] = color.TagBadge(t.Name, t.Color)
			}
			fmt.Printf("Tags:    %s\n", strings.Join(badges, " "))
		}
		if show("attendees") && len(result.Attendees) > 0 {
			emails := make([]string, len(result.Attendees))
			for i, a := range result.Attendees {
				emails[i] = a.Email
			}
			fmt.Printf("Attendees: %s\n", strings.Join(emails, ", "))
		}
		if show("summary") && result.Summary != nil {
			fmt.Printf("\n--- Summary ---\n%s\n", *result.Summary)
		}
		if show("transcript") && len(result.Transcript) > 0 {
			fmt.Printf("\n--- Transcript (%d entries) ---\n", len(result.Transcript))
			for _, t := range result.Transcript {
				fmt.Printf("[%s] %s\n", t.Role, t.Text)
			}
		}

		return nil
	},
}

func init() {
	sessionsGetCmd.Flags().String("fields", "", "Comma-separated sections to show (e.g. title,summary,transcript)")
	sessionsGetCmd.Flags().String("no-fields", "", "Comma-separated sections to hide (e.g. transcript,attendees)")
	sessionsCmd.AddCommand(sessionsGetCmd)
}

func resolveGetSections(fields, noFields string) map[string]bool {
	sections := map[string]bool{}

	if fields != "" {
		for _, f := range strings.Split(fields, ",") {
			sections[strings.TrimSpace(f)] = true
		}
		return sections
	}

	for _, s := range allGetSections {
		sections[s] = true
	}
	if noFields != "" {
		for _, f := range strings.Split(noFields, ",") {
			delete(sections, strings.TrimSpace(f))
		}
	}
	return sections
}
