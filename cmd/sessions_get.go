package cmd

import (
	"fmt"
	"strings"

	"github.com/cluely/cli/internal/api"
	"github.com/spf13/cobra"
)

var sessionsGetCmd = &cobra.Command{
	Use:   "get <session-id>",
	Short: "Get session details",
	Long:  "Display details and transcript for a specific session.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		jsonOut, _ := cmd.Flags().GetBool("json")

		input := map[string]string{"id": sessionID}

		if jsonOut {
			raw, err := api.CallRaw("sessions/get", input)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		}

		var result struct {
			ID        string  `json:"id"`
			State     string  `json:"state"`
			Title     *string `json:"title"`
			Summary   *string `json:"summary"`
			CreatedAt string  `json:"createdAt"`
			EndedAt   *string `json:"endedAt"`
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

		title := "(untitled)"
		if result.Title != nil {
			title = *result.Title
		}

		fmt.Printf("Session: %s\n", result.ID)
		fmt.Printf("Title:   %s\n", title)
		fmt.Printf("State:   %s\n", result.State)
		fmt.Printf("Created: %s\n", formatTime(result.CreatedAt))
		if result.EndedAt != nil {
			fmt.Printf("Ended:   %s\n", formatTime(*result.EndedAt))
		}

		if len(result.Attendees) > 0 {
			emails := make([]string, len(result.Attendees))
			for i, a := range result.Attendees {
				emails[i] = a.Email
			}
			fmt.Printf("Attendees: %s\n", strings.Join(emails, ", "))
		}

		if result.Summary != nil {
			fmt.Printf("\n--- Summary ---\n%s\n", *result.Summary)
		}

		if len(result.Transcript) > 0 {
			fmt.Printf("\n--- Transcript (%d entries) ---\n", len(result.Transcript))
			for _, t := range result.Transcript {
				fmt.Printf("[%s] %s\n", t.Role, t.Text)
			}
		}

		return nil
	},
}

func init() {
	sessionsCmd.AddCommand(sessionsGetCmd)
}
