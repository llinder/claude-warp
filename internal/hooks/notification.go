package hooks

import (
	"encoding/json"
	"os"

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
		return notify.Send("Claude Code", "Input needed")
	}

	msg := input.Message
	if msg == "" {
		msg = "Input needed"
	}

	return notify.Send("Claude Code", msg)
}
