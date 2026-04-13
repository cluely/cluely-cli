package color

import (
	"fmt"
	"strconv"
	"strings"
)

// TagBadge returns a terminal-formatted tag with colored background.
// Hex should be like "#ff5733" or "ff5733".
func TagBadge(name, hex string) string {
	r, g, b, ok := parseHex(hex)
	if !ok {
		return fmt.Sprintf("[%s]", name)
	}

	// Pick white or black foreground based on luminance
	fg := "\033[38;2;255;255;255m" // white
	if luminance(r, g, b) > 0.5 {
		fg = "\033[38;2;0;0;0m" // black
	}

	bg := fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
	reset := "\033[0m"

	return fmt.Sprintf("%s%s %s %s", bg, fg, name, reset)
}

func parseHex(hex string) (r, g, b uint8, ok bool) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, false
	}
	ri, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}
	gi, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}
	bi, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}
	return uint8(ri), uint8(gi), uint8(bi), true
}

// luminance returns relative luminance (0 = dark, 1 = bright).
func luminance(r, g, b uint8) float64 {
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255
}
