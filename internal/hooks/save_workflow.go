package hooks

import (
	"fmt"
	"strings"

	"github.com/llinder/claude-warp/internal/workflow"
)

// SaveWorkflowOpts holds parsed options for the save-workflow command.
type SaveWorkflowOpts struct {
	Name        string
	Command     string
	Description string
	RepoScoped  bool
	ArgSpecs    []string
	Tags        []string
}

// SaveWorkflow saves a command as a Warp workflow.
func SaveWorkflow(opts SaveWorkflowOpts) error {
	w := workflow.Workflow{
		Name:        opts.Name,
		Command:     opts.Command,
		Description: opts.Description,
		Tags:        opts.Tags,
	}

	for _, spec := range opts.ArgSpecs {
		arg, err := parseArgSpec(spec)
		if err != nil {
			return fmt.Errorf("invalid --arg %q: %w", spec, err)
		}
		w.Arguments = append(w.Arguments, arg)
	}

	path, err := workflow.Save(w, opts.RepoScoped)
	if err != nil {
		return err
	}

	scope := "personal"
	if opts.RepoScoped {
		scope = "repo"
	}
	fmt.Printf("Saved %s workflow to %s\n", scope, path)
	fmt.Println("Access it in Warp via Ctrl+Shift+R (Command Search)")
	return nil
}

// parseArgSpec parses "name:description:default" format.
// Default is optional. Use ~ for no default.
func parseArgSpec(spec string) (workflow.Argument, error) {
	parts := strings.SplitN(spec, ":", 3)
	if len(parts) < 1 || parts[0] == "" {
		return workflow.Argument{}, fmt.Errorf("argument name is required")
	}

	arg := workflow.Argument{
		Name:         parts[0],
		DefaultValue: nil, // nil marshals to ~ in YAML (no default)
	}

	if len(parts) >= 2 {
		arg.Description = parts[1]
	}
	if len(parts) >= 3 && parts[2] != "~" {
		arg.DefaultValue = parts[2]
	}

	return arg, nil
}
