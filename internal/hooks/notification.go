package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/llinder/claude-warp/internal/notify"
)

type notificationInput struct {
	Message string `json:"message"`
}

// Notification handles the Notification hook event.
// Forwards Claude Code notifications to Warp native notifications.
func Notification() error {
	var input notificationInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		return notify.Send(notifyTitle(), "Input needed")
	}

	msg := input.Message
	if msg == "" {
		msg = "Input needed"
	}

	// Re-set the tab title so it persists if Warp auto-title resets it
	notify.SetTabTitle(notifyTitle())

	return notify.Send(notifyTitle(), msg)
}

// notifyTitle returns a title that includes the current project directory
// so users can identify which tab/session needs attention.
func notifyTitle() string {
	dir, err := os.Getwd()
	if err != nil {
		return "Claude Code"
	}
	return "Claude Code - " + filepath.Base(dir)
}
