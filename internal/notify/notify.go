package notify

import (
	"fmt"
	"os"
	"strings"
)

// Send sends a Warp notification using OSC 777 escape sequences.
// Title and body are sanitized to prevent terminal escape injection.
func Send(title, body string) error {
	title = sanitize(title)
	body = sanitize(body)

	tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		return nil // silently fail if no tty (e.g. running in CI)
	}
	defer tty.Close()

	// OSC 777 format: \033]777;notify;<title>;<body>\007
	_, err = fmt.Fprintf(tty, "\033]777;notify;%s;%s\007", title, body)
	return err
}

// sanitize strips control characters (0x00-0x1F, 0x7F) to prevent
// terminal escape sequence injection.
func sanitize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r >= 0x20 && r != 0x7F {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Truncate truncates a string to maxLen characters, appending "..." if truncated.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
