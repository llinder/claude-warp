package hooks

import (
	"os"
	"path/filepath"

	"github.com/llinder/claude-warp/internal/notify"
)

// SessionStart handles the SessionStart hook event.
// Sets the Warp tab title for session identification.
func SessionStart() error {
	inWarp := os.Getenv("TERM_PROGRAM") == "WarpTerminal"

	if !inWarp {
		return nil
	}

	// Set Warp tab title to project name for easy identification
	if dir, err := os.Getwd(); err == nil {
		notify.SetTabTitle("Claude: " + filepath.Base(dir))
	}

	return nil
}

