package notify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Send sends a Warp notification using OSC 777 escape sequences.
// Skips the notification if Warp is the frontmost app (user is already looking at it).
// Title and body are sanitized to prevent terminal escape injection.
func Send(title, body string) error {
	if isWarpFocused() {
		return nil
	}

	title = sanitize(title)
	body = sanitize(body)

	tty, err := openTTY()
	if err != nil {
		return nil // silently fail if no tty (e.g. running in CI)
	}
	defer tty.Close()

	// OSC 777 format: \033]777;notify;<title>;<body>\007
	_, err = fmt.Fprintf(tty, "\033]777;notify;%s;%s\007", title, body)
	return err
}

// isWarpFocused checks if Warp is the frontmost application on macOS.
func isWarpFocused() bool {
	if runtime.GOOS != "darwin" {
		return false
	}
	out, err := exec.Command("osascript", "-e",
		`tell application "System Events" to get bundle identifier of first application process whose frontmost is true`).Output()
	if err != nil {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(string(out)), "dev.warp.Warp")
}

// openTTY tries /dev/tty first, then falls back to known TTY env vars.
// Claude Code subprocesses may not have a controlling terminal, so we
// check for the actual PTY path via environment variables.
func openTTY() (*os.File, error) {
	f, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err == nil {
		return f, nil
	}

	for _, env := range []string{"_P9K_SSH_TTY", "_P9K_TTY", "GPG_TTY"} {
		if path := os.Getenv(env); path != "" {
			if f, err := os.OpenFile(path, os.O_WRONLY, 0); err == nil {
				return f, nil
			}
		}
	}

	return nil, fmt.Errorf("no tty available")
}

// SetTabTitle sets the Warp tab title using OSC 0 escape sequences.
func SetTabTitle(title string) error {
	title = sanitize(title)

	tty, err := openTTY()
	if err != nil {
		return nil
	}
	defer tty.Close()

	// OSC 0 format: \033]0;<title>\007
	_, err = fmt.Fprintf(tty, "\033]0;%s\007", title)
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
