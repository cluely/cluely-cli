package version

import (
	"fmt"
	"strings"
	"testing"
)

func TestFull(t *testing.T) {
	t.Run("returns formatted version string", func(t *testing.T) {
		Version = "1.2.3"
		Commit = "abc1234"
		Date = "2024-01-15"

		got := Full()
		want := fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date)
		if got != want {
			t.Errorf("Full() = %q, want %q", got, want)
		}
	})

	t.Run("uses default values when not set", func(t *testing.T) {
		Version = "dev"
		Commit = "unknown"
		Date = "unknown"

		got := Full()
		want := "dev (commit unknown, built unknown)"
		if got != want {
			t.Errorf("Full() = %q, want %q", got, want)
		}
	})

	t.Run("contains version, commit, and date fields", func(t *testing.T) {
		Version = "2.0.0"
		Commit = "deadbeef"
		Date = "2025-06-01"

		got := Full()
		for _, substr := range []string{Version, Commit, Date} {
			if !strings.Contains(got, substr) {
				t.Errorf("Full() = %q, expected to contain %q", got, substr)
			}
		}
	})
}
