package cmd

import (
	"testing"
)

func strPtr(s string) *string { return &s }

// --- titleOrEmpty ---

func TestTitleOrEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  string
	}{
		{
			name:  "nil returns empty string",
			input: nil,
			want:  "",
		},
		{
			name:  "empty string returns empty string",
			input: strPtr(""),
			want:  "",
		},
		{
			name:  "whitespace-only string is trimmed to empty",
			input: strPtr("   "),
			want:  "",
		},
		{
			name:  "normal title is returned as-is",
			input: strPtr("My Session"),
			want:  "My Session",
		},
		{
			name:  "title with leading and trailing spaces is trimmed",
			input: strPtr("  trimmed  "),
			want:  "trimmed",
		},
		{
			name:  "tab-padded title is trimmed",
			input: strPtr("\tmy-title\t"),
			want:  "my-title",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := titleOrEmpty(tc.input)
			if got != tc.want {
				t.Errorf("titleOrEmpty() = %q, want %q", got, tc.want)
			}
		})
	}
}

// --- titleDisplay ---

func TestTitleDisplay(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  string
	}{
		{
			name:  "nil returns (untitled)",
			input: nil,
			want:  "(untitled)",
		},
		{
			name:  "empty string returns (untitled)",
			input: strPtr(""),
			want:  "(untitled)",
		},
		{
			name:  "whitespace-only string returns (untitled)",
			input: strPtr("   "),
			want:  "(untitled)",
		},
		{
			name:  "normal title is returned trimmed",
			input: strPtr("My Session"),
			want:  "My Session",
		},
		{
			name:  "title with surrounding spaces is trimmed",
			input: strPtr("  padded title  "),
			want:  "padded title",
		},
		{
			name:  "single character is not (untitled)",
			input: strPtr("X"),
			want:  "X",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := titleDisplay(tc.input)
			if got != tc.want {
				t.Errorf("titleDisplay() = %q, want %q", got, tc.want)
			}
		})
	}
}
