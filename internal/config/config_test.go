package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows path logic tested separately")
	}

	t.Run("uses XDG_CONFIG_HOME when set", func(t *testing.T) {
		xdg := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", xdg)

		got := Dir()
		want := filepath.Join(xdg, "cluely")
		if got != want {
			t.Errorf("Dir() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to ~/.config/cluely when XDG_CONFIG_HOME is unset", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")

		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("cannot determine home dir:", err)
		}

		got := Dir()
		want := filepath.Join(home, ".config", "cluely")
		if got != want {
			t.Errorf("Dir() = %q, want %q", got, want)
		}
	})

	t.Run("result always ends with 'cluely'", func(t *testing.T) {
		got := Dir()
		if !strings.HasSuffix(got, "cluely") {
			t.Errorf("Dir() = %q, expected path to end with 'cluely'", got)
		}
	})

	t.Run("result is an absolute path", func(t *testing.T) {
		got := Dir()
		if !filepath.IsAbs(got) {
			t.Errorf("Dir() = %q, expected absolute path", got)
		}
	})
}
