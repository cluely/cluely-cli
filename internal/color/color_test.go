package color

import (
	"strings"
	"testing"
)

// --- parseHex ---

func TestParseHex(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantR  uint8
		wantG  uint8
		wantB  uint8
		wantOK bool
	}{
		{
			name:   "valid hex with hash prefix",
			input:  "#ff5733",
			wantR:  0xff, wantG: 0x57, wantB: 0x33,
			wantOK: true,
		},
		{
			name:   "valid hex without hash prefix",
			input:  "ff5733",
			wantR:  0xff, wantG: 0x57, wantB: 0x33,
			wantOK: true,
		},
		{
			name:   "black",
			input:  "#000000",
			wantR:  0, wantG: 0, wantB: 0,
			wantOK: true,
		},
		{
			name:   "white",
			input:  "#ffffff",
			wantR:  255, wantG: 255, wantB: 255,
			wantOK: true,
		},
		{
			name:   "uppercase hex",
			input:  "#FF5733",
			wantR:  0xff, wantG: 0x57, wantB: 0x33,
			wantOK: true,
		},
		{
			name:   "too short",
			input:  "#fff",
			wantOK: false,
		},
		{
			name:   "too long",
			input:  "#ffffffff",
			wantOK: false,
		},
		{
			name:   "empty string",
			input:  "",
			wantOK: false,
		},
		{
			name:   "invalid hex characters",
			input:  "#zzzzzz",
			wantOK: false,
		},
		{
			name:   "invalid red component",
			input:  "#zz0000",
			wantOK: false,
		},
		{
			name:   "invalid green component",
			input:  "#00zz00",
			wantOK: false,
		},
		{
			name:   "invalid blue component",
			input:  "#0000zz",
			wantOK: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, g, b, ok := parseHex(tc.input)
			if ok != tc.wantOK {
				t.Errorf("parseHex(%q) ok = %v, want %v", tc.input, ok, tc.wantOK)
				return
			}
			if ok {
				if r != tc.wantR || g != tc.wantG || b != tc.wantB {
					t.Errorf("parseHex(%q) = (%d,%d,%d), want (%d,%d,%d)",
						tc.input, r, g, b, tc.wantR, tc.wantG, tc.wantB)
				}
			}
		})
	}
}

// --- luminance ---

func TestLuminance(t *testing.T) {
	tests := []struct {
		name string
		r, g, b uint8
		wantMin float64
		wantMax float64
	}{
		{
			name:    "black is 0",
			r: 0, g: 0, b: 0,
			wantMin: 0, wantMax: 0,
		},
		{
			name:    "white is 1",
			r: 255, g: 255, b: 255,
			wantMin: 1, wantMax: 1,
		},
		{
			name:    "pure red is roughly 0.299",
			r: 255, g: 0, b: 0,
			wantMin: 0.29, wantMax: 0.31,
		},
		{
			name:    "pure green is roughly 0.587",
			r: 0, g: 255, b: 0,
			wantMin: 0.58, wantMax: 0.60,
		},
		{
			name:    "pure blue is roughly 0.114",
			r: 0, g: 0, b: 255,
			wantMin: 0.11, wantMax: 0.12,
		},
		{
			name:    "mid-grey is around 0.5",
			r: 128, g: 128, b: 128,
			wantMin: 0.49, wantMax: 0.51,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := luminance(tc.r, tc.g, tc.b)
			if got < tc.wantMin || got > tc.wantMax {
				t.Errorf("luminance(%d,%d,%d) = %f, want [%f, %f]",
					tc.r, tc.g, tc.b, got, tc.wantMin, tc.wantMax)
			}
		})
	}
}

// --- TagBadge ---

func TestTagBadge(t *testing.T) {
	t.Run("invalid hex returns fallback", func(t *testing.T) {
		got := TagBadge("mytag", "notahex")
		want := "[mytag]"
		if got != want {
			t.Errorf("TagBadge(invalid hex) = %q, want %q", got, want)
		}
	})

	t.Run("empty hex returns fallback", func(t *testing.T) {
		got := TagBadge("label", "")
		want := "[label]"
		if got != want {
			t.Errorf("TagBadge(empty hex) = %q, want %q", got, want)
		}
	})

	t.Run("dark background uses white foreground", func(t *testing.T) {
		// #000000 is black — luminance 0, so white fg is expected
		got := TagBadge("dark", "#000000")
		// White fg ANSI code
		if !strings.Contains(got, "255;255;255") {
			t.Errorf("dark background should use white fg, got: %q", got)
		}
		if !strings.Contains(got, "dark") {
			t.Errorf("TagBadge result should contain tag name, got: %q", got)
		}
	})

	t.Run("light background uses black foreground", func(t *testing.T) {
		// #ffffff is white — luminance 1 (>0.5), so black fg is expected
		got := TagBadge("light", "#ffffff")
		// Black fg ANSI code
		if !strings.Contains(got, "0;0;0") {
			t.Errorf("light background should use black fg, got: %q", got)
		}
		if !strings.Contains(got, "light") {
			t.Errorf("TagBadge result should contain tag name, got: %q", got)
		}
	})

	t.Run("background ANSI code embedded in output", func(t *testing.T) {
		// #ff0000 → 255,0,0
		got := TagBadge("red", "#ff0000")
		if !strings.Contains(got, "48;2;255;0;0") {
			t.Errorf("expected bg ANSI code for red, got: %q", got)
		}
	})

	t.Run("output contains reset code", func(t *testing.T) {
		got := TagBadge("tag", "#123456")
		if !strings.Contains(got, "\033[0m") {
			t.Errorf("expected ANSI reset code in output, got: %q", got)
		}
	})

	t.Run("hex without hash prefix is accepted", func(t *testing.T) {
		withHash := TagBadge("t", "#aabbcc")
		withoutHash := TagBadge("t", "aabbcc")
		if withHash != withoutHash {
			t.Errorf("hex with/without # should produce same result: %q vs %q", withHash, withoutHash)
		}
	})
}
