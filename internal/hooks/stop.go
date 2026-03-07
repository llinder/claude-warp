package hooks

import (
	"encoding/json"
	"os"

	"github.com/llinder/claude-warp/internal/notify"
	"github.com/llinder/claude-warp/internal/transcript"
)

// stopInput is the JSON input provided by Claude Code on the Stop hook.
type stopInput struct {
	TranscriptPath      string `json:"transcript_path"`
	LastAssistantMessage string `json:"last_assistant_message"`
}

// Stop handles the Stop hook event.
func Stop() error {
	var input stopInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		return notify.Send(notifyTitle(), "Done")
	}

	msg := buildStopMessage(input)
	return notify.Send(notifyTitle(), msg)
}

func buildStopMessage(input stopInput) string {
	// Prefer the last assistant message from the hook input
	if input.LastAssistantMessage != "" {
		return notify.Truncate(input.LastAssistantMessage, 80)
	}

	// Fall back to transcript parsing
	if input.TranscriptPath != "" {
		messages, err := transcript.Parse(input.TranscriptPath)
		if err == nil {
			if resp := transcript.ExtractLastResponse(messages); resp != "" {
				return notify.Truncate(resp, 80)
			}
		}
	}

	return "Done"
}
