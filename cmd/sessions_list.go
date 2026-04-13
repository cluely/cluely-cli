package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cluely/cli/internal/api"
	"github.com/spf13/cobra"
)

var sessionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sessions",
	Long:  "List your meeting sessions, newest first.",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")
		state, _ := cmd.Flags().GetString("state")
		since, _ := cmd.Flags().GetString("since")
		jsonOut, _ := cmd.Flags().GetBool("json")

		input := map[string]interface{}{
			"limit": limit,
		}
		if cursor != "" {
			input["cursor"] = cursor
		}
		if state != "" {
			input["state"] = state
		}
		if since != "" {
			d, err := parseDuration(since)
			if err != nil {
				return fmt.Errorf("invalid --since value %q: %w (examples: 24h, 7d, 30m)", since, err)
			}
			input["createdAfter"] = time.Now().Add(-d).UTC().Format(time.RFC3339)
		}

		if jsonOut {
			raw, err := api.CallRaw("sessions/list", input)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		}

		var result struct {
			Items []struct {
				ID        string  `json:"id"`
				State     string  `json:"state"`
				Title     *string `json:"title"`
				CreatedAt string  `json:"createdAt"`
				EndedAt   *string `json:"endedAt"`
			} `json:"items"`
			NextCursor *string `json:"nextCursor"`
			Total      int     `json:"total"`
		}

		if err := api.Call("sessions/list", input, &result); err != nil {
			return err
		}

		if len(result.Items) == 0 {
			fmt.Println("No sessions found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tSTATE\tTITLE\tCREATED\n")

		for _, s := range result.Items {
			title := "-"
			if s.Title != nil {
				title = truncate(*s.Title, 50)
			}
			created := formatTime(s.CreatedAt)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.State, title, created)
		}
		w.Flush()

		fmt.Printf("\nShowing %d of %d sessions.", len(result.Items), result.Total)
		if result.NextCursor != nil {
			fmt.Printf(" Use --cursor %q to see more.", *result.NextCursor)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	sessionsListCmd.Flags().IntP("limit", "n", 20, "Number of sessions to show")
	sessionsListCmd.Flags().String("cursor", "", "Pagination cursor for next page")
	sessionsListCmd.Flags().String("state", "", "Filter by state (ongoing, analyzing, finished)")
	sessionsListCmd.Flags().String("since", "", "Show sessions created in the last duration (e.g. 24h, 7d, 30m)")
	sessionsCmd.AddCommand(sessionsListCmd)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatTime(raw string) string {
	t, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return raw
	}
	return t.Local().Format("Jan 02, 15:04")
}

func prettyJSON(raw json.RawMessage) string {
	var buf []byte
	buf, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return string(raw)
	}
	return string(buf)
}

// parseDuration parses a human-friendly duration string like "24h", "7d", "30m".
// Supports: m (minutes), h (hours), d (days).
func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return 0, fmt.Errorf("too short")
	}

	unit := s[len(s)-1]
	numStr := s[:len(s)-1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", numStr)
	}

	switch unit {
	case 'm':
		return time.Duration(num) * time.Minute, nil
	case 'h':
		return time.Duration(num) * time.Hour, nil
	case 'd':
		return time.Duration(num) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown unit %q, use m/h/d", string(unit))
	}
}
