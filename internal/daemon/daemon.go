package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cluely/cli/internal/config"
)

const serviceName = "com.cluely.watch"

// Install creates and starts the background watch service.
// The service runs: cluely sessions watch --exec "<execCmd>"
func Install(execCmd string) error {
	cluely, err := findBinary()
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "darwin":
		return installLaunchd(cluely, execCmd)
	case "linux":
		return installSystemd(cluely, execCmd)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// Uninstall stops and removes the background watch service.
func Uninstall() error {
	switch runtime.GOOS {
	case "darwin":
		return uninstallLaunchd()
	case "linux":
		return uninstallSystemd()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// Status returns whether the service is running and the exec command if configured.
func Status() (running bool, execCmd string, err error) {
	switch runtime.GOOS {
	case "darwin":
		return statusLaunchd()
	case "linux":
		return statusSystemd()
	default:
		return false, "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// LogPath returns the path to the service log file.
func LogPath() string {
	return filepath.Join(logDir(), "watch.log")
}

func logDir() string {
	dir := filepath.Join(config.Dir(), "logs")
	os.MkdirAll(dir, 0o755)
	return dir
}

func findBinary() (string, error) {
	path, err := exec.LookPath("cluely")
	if err != nil {
		// Fall back to the current executable
		path, err = os.Executable()
		if err != nil {
			return "", fmt.Errorf("cannot find cluely binary: %w", err)
		}
	}
	return filepath.Abs(path)
}

// --- macOS launchd ---

func launchdPlistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", serviceName+".plist")
}

func installLaunchd(binary, execCmd string) error {
	logFile := LogPath()

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>sessions</string>
        <string>watch</string>
        <string>--exec</string>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>%s</string>
    <key>StandardErrorPath</key>
    <string>%s</string>
</dict>
</plist>`, serviceName, binary, escapeXML(execCmd), logFile, logFile)

	path := launchdPlistPath()
	if err := os.WriteFile(path, []byte(plist), 0o644); err != nil {
		return fmt.Errorf("write plist: %w", err)
	}

	// Unload first in case it's already loaded (ignore errors)
	exec.Command("launchctl", "unload", path).Run()

	if out, err := exec.Command("launchctl", "load", path).CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl load: %s", strings.TrimSpace(string(out)))
	}

	return nil
}

func uninstallLaunchd() error {
	path := launchdPlistPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("service is not installed")
	}

	if out, err := exec.Command("launchctl", "unload", path).CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl unload: %s", strings.TrimSpace(string(out)))
	}

	return os.Remove(path)
}

func statusLaunchd() (bool, string, error) {
	path := launchdPlistPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, "", nil
	}

	out, err := exec.Command("launchctl", "list", serviceName).CombinedOutput()
	if err != nil {
		return false, "", nil // Not loaded
	}

	// Parse the exec command from the plist
	plistBytes, _ := os.ReadFile(path)
	execCmd := extractExecFromPlist(string(plistBytes))

	return strings.Contains(string(out), serviceName), execCmd, nil
}

func extractExecFromPlist(plist string) string {
	// Find the string after --exec in ProgramArguments
	marker := "<string>--exec</string>"
	idx := strings.Index(plist, marker)
	if idx == -1 {
		return ""
	}
	rest := plist[idx+len(marker):]
	start := strings.Index(rest, "<string>")
	end := strings.Index(rest, "</string>")
	if start == -1 || end == -1 || end < start+len("<string>") {
		return ""
	}
	return rest[start+len("<string>") : end]
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// --- Linux systemd ---

func systemdServicePath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".config", "systemd", "user")
	os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "cluely-watch.service")
}

func installSystemd(binary, execCmd string) error {
	logFile := LogPath()

	unit := fmt.Sprintf(`[Unit]
Description=Cluely Session Watcher
After=network-online.target

[Service]
Type=simple
ExecStart=%s sessions watch --exec "%s"
Restart=always
RestartSec=5
StandardOutput=append:%s
StandardError=append:%s

[Install]
WantedBy=default.target
`, binary, strings.ReplaceAll(execCmd, `"`, `\"`), logFile, logFile)

	path := systemdServicePath()
	if err := os.WriteFile(path, []byte(unit), 0o644); err != nil {
		return fmt.Errorf("write service file: %w", err)
	}

	cmds := [][]string{
		{"systemctl", "--user", "daemon-reload"},
		{"systemctl", "--user", "enable", "cluely-watch.service"},
		{"systemctl", "--user", "start", "cluely-watch.service"},
	}
	for _, args := range cmds {
		if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
			return fmt.Errorf("%s: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
		}
	}

	return nil
}

func uninstallSystemd() error {
	path := systemdServicePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("service is not installed")
	}

	cmds := [][]string{
		{"systemctl", "--user", "stop", "cluely-watch.service"},
		{"systemctl", "--user", "disable", "cluely-watch.service"},
	}
	for _, args := range cmds {
		exec.Command(args[0], args[1:]...).Run() // Ignore errors
	}

	exec.Command("systemctl", "--user", "daemon-reload").Run()

	return os.Remove(path)
}

func statusSystemd() (bool, string, error) {
	path := systemdServicePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, "", nil
	}

	out, err := exec.Command("systemctl", "--user", "is-active", "cluely-watch.service").CombinedOutput()
	running := err == nil && strings.TrimSpace(string(out)) == "active"

	unitBytes, _ := os.ReadFile(path)
	execCmd := extractExecFromUnit(string(unitBytes))

	return running, execCmd, nil
}

func extractExecFromUnit(unit string) string {
	for _, line := range strings.Split(unit, "\n") {
		if strings.HasPrefix(line, "ExecStart=") {
			// ExecStart=/path/to/cluely sessions watch --exec "cmd"
			parts := strings.SplitN(line, "--exec ", 2)
			if len(parts) == 2 {
				return strings.Trim(strings.TrimSpace(parts[1]), `"`)
			}
		}
	}
	return ""
}
