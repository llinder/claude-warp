package hooks

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/llinder/claude-warp/internal/notify"
	"github.com/llinder/claude-warp/internal/transcript"
)

// stopInput is the JSON input provided by Claude Code on the Stop hook.
type stopInput struct {
	TranscriptPath string `json:"transcript_path"`
}

// Stop handles the Stop hook event.
// Parses the session transcript and sends a rich notification with the task summary.
func Stop() error {
	var input stopInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		// Fall back to simple notification if we can't read input
		return notify.Send(notifyTitle(), "Task completed")
	}

	msg := buildStopMessage(input.TranscriptPath)
	return notify.Send(notifyTitle(), msg)
}

func buildStopMessage(transcriptPath string) string {
	if transcriptPath == "" {
		return "Task completed"
	}

	messages, err := transcript.Parse(transcriptPath)
	if err != nil || len(messages) == 0 {
		return "Task completed"
	}

	prompt := transcript.ExtractFirstPrompt(messages)
	response := transcript.ExtractLastResponse(messages)
	filesChanged := transcript.CountFilesChanged(messages)
	bashCmds := transcript.ExtractBashCommands(messages)

	var msg string

	// Start with the prompt if available
	if prompt != "" {
		msg = fmt.Sprintf("\"%s\"", notify.Truncate(prompt, 50))
	}

	// Add summary stats
	var stats []string
	if filesChanged > 0 {
		stats = append(stats, fmt.Sprintf("%d files changed", filesChanged))
	}
	if len(bashCmds) > 0 {
		stats = append(stats, fmt.Sprintf("%d commands run", len(bashCmds)))
	}

	if msg != "" && len(stats) > 0 {
		msg += " | "
	}
	for i, s := range stats {
		if i > 0 {
			msg += ", "
		}
		msg += s
	}

	// Add response snippet
	if response != "" {
		resp := notify.Truncate(response, 100)
		if msg != "" {
			msg += " -> " + resp
		} else {
			msg = resp
		}
	}

	if msg == "" {
		msg = "Task completed"
	}

	return notify.Truncate(msg, 200)
}
