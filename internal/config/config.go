package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	APIURL = "https://api.v2.cluely.com"
	WebURL = "https://v2.cluely.com"
)

// Dir returns the platform-appropriate config directory for cluely.
func Dir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "cluely")
	}
	home, _ := os.UserHomeDir()
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "cluely")
	}
	return filepath.Join(home, ".config", "cluely")
}
