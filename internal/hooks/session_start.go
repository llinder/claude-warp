package hooks

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/llinder/claude-warp/internal/workflow"
)

// SessionStart handles the SessionStart hook event.
// Detects Warp, discovers project context, and injects a system message
// teaching Claude about available Warp integration features.
func SessionStart() error {
	inWarp := os.Getenv("TERM_PROGRAM") == "WarpTerminal"

	if !inWarp {
		resp := map[string]string{
			"systemMessage": "Warp plugin installed but not running in Warp terminal. " +
				"Install Warp (https://warp.dev) for native notifications and workflow integration.",
		}
		return json.NewEncoder(os.Stdout).Encode(resp)
	}

	var parts []string
	parts = append(parts, "Warp terminal detected. Integration features available:")
	parts = append(parts, "")

	// Notifications
	parts = append(parts, "NOTIFICATIONS: Task completion and input-needed notifications are sent automatically via Warp native notifications.")

	// Discover existing workflows
	personalWf, _ := workflow.List(false)
	repoWf, _ := workflow.List(true)

	if len(personalWf) > 0 || len(repoWf) > 0 {
		parts = append(parts, "")
		parts = append(parts, "EXISTING WARP WORKFLOWS:")
		if len(repoWf) > 0 {
			parts = append(parts, fmt.Sprintf("  Repo workflows (%s): %d found", workflow.RepoDir(), len(repoWf)))
			for _, w := range repoWf {
				parts = append(parts, fmt.Sprintf("    - %s: %s", w.Name, w.Description))
			}
		}
		if len(personalWf) > 0 {
			parts = append(parts, fmt.Sprintf("  Personal workflows (%s): %d found", workflow.PersonalDir(), len(personalWf)))
		}
	}

	// CLI instructions for Claude
	bin := selfPath()
	parts = append(parts, "")
	parts = append(parts, "WARP WORKFLOW COMMANDS: You can save useful commands as Warp workflows that the user can access via Ctrl+Shift+R.")
	parts = append(parts, fmt.Sprintf("  Save personal workflow:    %s save-workflow --name <name> --command <cmd> [--description <desc>] [--arg <name:description:default>]...", bin))
	parts = append(parts, fmt.Sprintf("  Save repo workflow:        %s save-workflow --repo --name <name> --command <cmd> [--description <desc>] [--arg <name:description:default>]...", bin))
	parts = append(parts, "  Arguments use {{arg_name}} syntax in the command string.")
	parts = append(parts, "")
	parts = append(parts, "WARP LAUNCH CONFIGS: You can create multi-tab/pane dev environment layouts.")
	parts = append(parts, fmt.Sprintf("  Save launch config:        %s save-launch --name <name> --tab <title:cwd:command> [--tab <title:cwd:command>]...", bin))
	parts = append(parts, "")
	parts = append(parts, "When you discover useful or repeated commands during a session, proactively suggest saving them as Warp workflows.")

	resp := map[string]string{
		"systemMessage": strings.Join(parts, "\n"),
	}
	return json.NewEncoder(os.Stdout).Encode(resp)
}

func selfPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "claude-warp"
	}
	return exe
}
