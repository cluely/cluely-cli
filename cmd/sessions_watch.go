package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cluely/cli/internal/api"
	"github.com/spf13/cobra"
)

type ongoingSession struct {
	ID    string  `json:"id"`
	Title *string `json:"title"`
}

var sessionsWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for session completions",
	Long: `Watch for sessions to finish and optionally run a command.

Runs continuously until interrupted with Ctrl+C. When a session finishes,
prints its info and runs the --exec command if provided.

The --exec command has access to these environment variables:
  CLUELY_SESSION_ID      Session ID
  CLUELY_SESSION_TITLE   Session title (if available)`,
	Example: `  cluely sessions watch
  cluely sessions watch --exec "echo \$CLUELY_SESSION_TITLE finished"
  cluely sessions watch --exec "cluely sessions get \$CLUELY_SESSION_ID --json | ./process.sh"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		execCmd, _ := cmd.Flags().GetString("exec")
		done := cmd.Context().Done()

		fmt.Println("Watching for session completions... (Ctrl+C to stop)")

		// Track sessions we're already watching to avoid duplicates
		watching := map[string]bool{}

		for {
			// Fetch ongoing sessions
			ongoing, err := fetchOngoingSessions()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to fetch sessions: %v\n", err)
				select {
				case <-done:
					return nil
				case <-time.After(10 * time.Second):
					continue
				}
			}

			// Start watching new sessions
			for _, s := range ongoing {
				if watching[s.ID] {
					continue
				}
				watching[s.ID] = true
				title := "(untitled)"
				if s.Title != nil {
					title = *s.Title
				}
				fmt.Printf("Watching session: %s — %s\n", s.ID, title)

				go func(session ongoingSession) {
					waitForSession(done, session, execCmd)
					delete(watching, session.ID)
				}(s)
			}

			// Wait before re-checking for new sessions
			select {
			case <-done:
				return nil
			case <-time.After(15 * time.Second):
			}
		}
	},
}

func init() {
	sessionsWatchCmd.Flags().String("exec", "", "Command to run when a session finishes")
	sessionsCmd.AddCommand(sessionsWatchCmd)
}

func fetchOngoingSessions() ([]ongoingSession, error) {
	var result struct {
		Items []ongoingSession `json:"items"`
	}
	input := map[string]interface{}{
		"limit": 50,
		"state": "ongoing",
	}
	if err := api.Call("sessions/list", input, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// waitForSession calls the waitFor endpoint in a loop until the session finishes.
func waitForSession(done <-chan struct{}, session ongoingSession, execCmd string) {
	for {
		select {
		case <-done:
			return
		default:
		}

		var result struct {
			Status string `json:"status"`
		}
		err := api.Call("sessions/waitFor", map[string]string{
			"id":    session.ID,
			"state": "finished",
		}, &result)

		if err != nil {
			// Session may have been deleted or we lost connectivity — stop watching it
			return
		}

		if result.Status == "fulfilled" {
			title := "(untitled)"
			if session.Title != nil {
				title = *session.Title
			}
			fmt.Printf("Session finished: %s — %s\n", session.ID, title)

			if execCmd != "" {
				runExecCommand(execCmd, session)
			}
			return
		}

		// timed-out: loop and try again
	}
}

func runExecCommand(command string, session ongoingSession) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"CLUELY_SESSION_ID="+session.ID,
		"CLUELY_SESSION_TITLE="+titleOrEmpty(session.Title),
	)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "exec error: %v\n", err)
	}
}

func titleOrEmpty(t *string) string {
	if t != nil {
		return strings.TrimSpace(*t)
	}
	return ""
}
