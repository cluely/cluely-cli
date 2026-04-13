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
	"github.com/cluely/cli/internal/color"
	"github.com/spf13/cobra"
)

type sessionTag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

var allColumns = []string{"id", "state", "title", "tags", "created"}

var sessionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sessions",
	Long:  "List your meeting sessions, newest first.",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")
		state, _ := cmd.Flags().GetString("state")
		since, _ := cmd.Flags().GetString("since")
		tag, _ := cmd.Flags().GetString("tag")
		jsonOut, _ := cmd.Flags().GetBool("json")
		fields, _ := cmd.Flags().GetString("fields")
		noFields, _ := cmd.Flags().GetString("no-fields")

		columns := resolveColumns(fields, noFields)

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
		if tag != "" {
			input["tagIds"] = []string{tag}
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
				ID        string       `json:"id"`
				State     string       `json:"state"`
				Title     *string      `json:"title"`
				Tags      []sessionTag `json:"tags"`
				CreatedAt string       `json:"createdAt"`
				EndedAt   *string      `json:"endedAt"`
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

		show := func(col string) bool { return columns[col] }

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		// Header
		var hdr []string
		if show("id") { hdr = append(hdr, "ID") }
		if show("state") { hdr = append(hdr, "STATE") }
		if show("title") { hdr = append(hdr, "TITLE") }
		if show("tags") { hdr = append(hdr, "TAGS") }
		if show("created") { hdr = append(hdr, "CREATED") }
		fmt.Fprintln(w, strings.Join(hdr, "\t"))

		// Rows
		for _, s := range result.Items {
			var row []string
			if show("id") { row = append(row, s.ID) }
			if show("state") { row = append(row, s.State) }
			if show("title") {
				t := "-"
				if s.Title != nil { t = truncate(*s.Title, 50) }
				row = append(row, t)
			}
			if show("tags") { row = append(row, formatTags(s.Tags)) }
			if show("created") { row = append(row, formatTime(s.CreatedAt)) }
			fmt.Fprintln(w, strings.Join(row, "\t"))
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
	sessionsListCmd.Flags().String("tag", "", "Filter by tag ID")
	sessionsListCmd.Flags().String("fields", "", "Comma-separated columns to show (e.g. id,title,tags)")
	sessionsListCmd.Flags().String("no-fields", "", "Comma-separated columns to hide (e.g. tags,state)")
	sessionsCmd.AddCommand(sessionsListCmd)
}

// resolveColumns returns the set of columns to display.
// --fields includes only those columns, --no-fields excludes them from the default set.
func resolveColumns(fields, noFields string) map[string]bool {
	cols := map[string]bool{}

	if fields != "" {
		for _, f := range strings.Split(fields, ",") {
			cols[strings.TrimSpace(f)] = true
		}
		return cols
	}

	for _, c := range allColumns {
		cols[c] = true
	}
	if noFields != "" {
		for _, f := range strings.Split(noFields, ",") {
			delete(cols, strings.TrimSpace(f))
		}
	}
	return cols
}

func formatTags(tags []sessionTag) string {
	if len(tags) == 0 {
		return "-"
	}
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = color.TagBadge(t.Name, t.Color)
	}
	return strings.Join(parts, " ")
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
