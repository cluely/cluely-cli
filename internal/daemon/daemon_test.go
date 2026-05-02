package daemon

import (
	"testing"
)

// --- extractExecFromPlist ---

func TestExtractExecFromPlist(t *testing.T) {
	tests := []struct {
		name  string
		plist string
		want  string
	}{
		{
			name: "normal plist with exec command",
			plist: `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/cluely</string>
        <string>sessions</string>
        <string>watch</string>
        <string>--exec</string>
        <string>echo hello</string>
    </array>
</dict>
</plist>`,
			want: "echo hello",
		},
		{
			name: "exec command with spaces and special chars",
			plist: `<string>--exec</string>
<string>./on-complete.sh arg1 arg2</string>`,
			want: "./on-complete.sh arg1 arg2",
		},
		{
			name:  "no --exec marker",
			plist: `<string>sessions</string><string>watch</string>`,
			want:  "",
		},
		{
			name:  "marker present but no following string tags",
			plist: `<string>--exec</string> nothing here`,
			want:  "",
		},
		{
			name:  "marker present, opening tag but no closing tag",
			plist: `<string>--exec</string><string>echo`,
			want:  "",
		},
		{
			name:  "closing tag appears before opening tag (malformed – bounds check returns empty)",
			plist: `<string>--exec</string></string><string>val</string>`,
			want:  "",
		},
		{
			name:  "empty exec value",
			plist: `<string>--exec</string><string></string>`,
			want:  "",
		},
		{
			name:  "empty plist",
			plist: "",
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractExecFromPlist(tc.plist)
			if got != tc.want {
				t.Errorf("extractExecFromPlist() = %q, want %q", got, tc.want)
			}
		})
	}
}

// --- extractExecFromUnit ---

func TestExtractExecFromUnit(t *testing.T) {
	tests := []struct {
		name string
		unit string
		want string
	}{
		{
			name: "normal systemd unit",
			unit: `[Unit]
Description=Cluely Session Watcher

[Service]
Type=simple
ExecStart=/usr/local/bin/cluely sessions watch --exec "echo done"
Restart=always
`,
			want: "echo done",
		},
		{
			name: "exec command without spaces",
			unit: `ExecStart=/usr/bin/cluely sessions watch --exec "script.sh"
`,
			want: "script.sh",
		},
		{
			name: "exec command with embedded spaces",
			unit: `ExecStart=/usr/bin/cluely sessions watch --exec "do thing arg1 arg2"
`,
			want: "do thing arg1 arg2",
		},
		{
			name: "no ExecStart line",
			unit: `[Unit]
Description=Test
`,
			want: "",
		},
		{
			name: "ExecStart without --exec",
			unit: `ExecStart=/usr/bin/cluely sessions watch
`,
			want: "",
		},
		{
			name:  "empty unit",
			unit:  "",
			want:  "",
		},
		{
			name: "multiple lines, ExecStart in the middle",
			unit: `[Service]
Type=simple
ExecStart=/path/cluely sessions watch --exec "my-hook.sh"
Restart=always
`,
			want: "my-hook.sh",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractExecFromUnit(tc.unit)
			if got != tc.want {
				t.Errorf("extractExecFromUnit() = %q, want %q", got, tc.want)
			}
		})
	}
}

// --- escapeXML ---

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no special characters",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "ampersand",
			input: "a&b",
			want:  "a&amp;b",
		},
		{
			name:  "less-than",
			input: "a<b",
			want:  "a&lt;b",
		},
		{
			name:  "greater-than",
			input: "a>b",
			want:  "a&gt;b",
		},
		{
			name:  "double-quote",
			input: `say "hello"`,
			want:  "say &quot;hello&quot;",
		},
		{
			name:  "all special characters combined",
			input: `<script>&alert("xss")</script>`,
			want:  "&lt;script&gt;&amp;alert(&quot;xss&quot;)&lt;/script&gt;",
		},
		{
			name:  "multiple ampersands",
			input: "a&b&c",
			want:  "a&amp;b&amp;c",
		},
		{
			name:  "already escaped string is double-escaped",
			input: "&amp;",
			want:  "&amp;amp;",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := escapeXML(tc.input)
			if got != tc.want {
				t.Errorf("escapeXML(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
