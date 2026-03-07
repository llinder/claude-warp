package hooks

import (
	"fmt"
	"strings"

	"github.com/llinder/claude-warp/internal/launch"
)

// SaveLaunchOpts holds parsed options for the save-launch command.
type SaveLaunchOpts struct {
	Name     string
	TabSpecs []string
}

// SaveLaunch creates a Warp launch configuration.
func SaveLaunch(opts SaveLaunchOpts) error {
	var tabs []launch.Tab
	for i, spec := range opts.TabSpecs {
		tab, err := parseTabSpec(spec)
		if err != nil {
			return fmt.Errorf("invalid --tab %q: %w", spec, err)
		}
		if i == 0 {
			tab.Layout.IsFocused = true
		}
		tabs = append(tabs, tab)
	}

	config := launch.Config{
		Name:              opts.Name,
		ActiveWindowIndex: 0,
		Windows: []launch.Window{
			{
				ActiveTabIndex: 0,
				Tabs:           tabs,
			},
		},
	}

	path, err := launch.Save(config)
	if err != nil {
		return err
	}

	fmt.Printf("Saved launch configuration to %s\n", path)
	fmt.Printf("Open it in Warp or use: open 'warp://launch/%s'\n", path)
	return nil
}

// parseTabSpec parses "title:cwd:command" format.
func parseTabSpec(spec string) (launch.Tab, error) {
	parts := strings.SplitN(spec, ":", 3)
	if len(parts) < 1 || parts[0] == "" {
		return launch.Tab{}, fmt.Errorf("tab title is required")
	}

	tab := launch.Tab{
		Title:  parts[0],
		Layout: launch.Layout{},
	}

	if len(parts) >= 2 && parts[1] != "" {
		tab.Layout.Cwd = parts[1]
	}
	if len(parts) >= 3 && parts[2] != "" {
		tab.Layout.Commands = []launch.Command{{Exec: parts[2]}}
	}

	return tab, nil
}
