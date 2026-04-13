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

type watchSession struct {
	ID    string  `json:"id"`
	Title *string `json:"title"`
}

var sessionsWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for session events",
	Long: `Watch for session starts and completions, optionally running a command for each event.

Runs continuously until interrupted with Ctrl+C.

The --exec command has access to these environment variables:
  CLUELY_EVENT           Event type: "start" or "end"
  CLUELY_SESSION_ID      Session ID
  CLUELY_SESSION_TITLE   Session title (if available)`,
	Example: `  cluely sessions watch
  cluely sessions watch --exec "echo \$CLUELY_EVENT: \$CLUELY_SESSION_TITLE"
  cluely sessions watch --on end --exec "./on-complete.sh"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		execCmd, _ := cmd.Flags().GetString("exec")
		onFilter, _ := cmd.Flags().GetString("on")
		done := cmd.Context().Done()

		if onFilter != "" && onFilter != "start" && onFilter != "end" {
			return fmt.Errorf("--on must be 'start' or 'end'")
		}

		fmt.Println("Watching for session events... (Ctrl+C to stop)")

		watching := map[string]bool{}

		for {
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

			for _, s := range ongoing {
				if watching[s.ID] {
					continue
				}
				watching[s.ID] = true

				fmt.Printf("Session started: %s — %s\n", s.ID, titleDisplay(s.Title))

				if onFilter == "" || onFilter == "start" {
					if execCmd != "" {
						runExecCommand(execCmd, s, "start")
					}
				}

				go func(session watchSession) {
					waitForSession(done, session, execCmd, onFilter)
					delete(watching, session.ID)
				}(s)
			}

			select {
			case <-done:
				return nil
			case <-time.After(15 * time.Second):
			}
		}
	},
}

func init() {
	sessionsWatchCmd.Flags().String("exec", "", "Command to run on session events")
	sessionsWatchCmd.Flags().String("on", "", "Filter events: 'start' or 'end' (default: both)")
	sessionsCmd.AddCommand(sessionsWatchCmd)
}

func fetchOngoingSessions() ([]watchSession, error) {
	var result struct {
		Items []watchSession `json:"items"`
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

func waitForSession(done <-chan struct{}, session watchSession, execCmd, onFilter string) {
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
			fmt.Fprintf(os.Stderr, "Warning: waitFor %s failed: %v\n", session.ID, err)
			return
		}

		if result.Status == "fulfilled" {
			fmt.Printf("Session ended: %s — %s\n", session.ID, titleDisplay(session.Title))

			if onFilter == "" || onFilter == "end" {
				if execCmd != "" {
					runExecCommand(execCmd, session, "end")
				}
			}
			return
		}
	}
}

func runExecCommand(command string, session watchSession, event string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"CLUELY_EVENT="+event,
		"CLUELY_SESSION_ID="+session.ID,
		"CLUELY_SESSION_TITLE="+titleOrEmpty(session.Title),
	)

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "exec error: %v\n", err)
		return
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "exec error: %v\n", err)
		}
	}()
}

func titleOrEmpty(t *string) string {
	if t != nil {
		return strings.TrimSpace(*t)
	}
	return ""
}

func titleDisplay(t *string) string {
	if t != nil && strings.TrimSpace(*t) != "" {
		return strings.TrimSpace(*t)
	}
	return "(untitled)"
}
