package launch

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents a Warp launch configuration.
type Config struct {
	Name              string   `yaml:"name"`
	ActiveWindowIndex int      `yaml:"active_window_index"`
	Windows           []Window `yaml:"windows"`
}

// Window represents a window in a launch configuration.
type Window struct {
	ActiveTabIndex int   `yaml:"active_tab_index"`
	Tabs           []Tab `yaml:"tabs"`
}

// Tab represents a tab in a window.
type Tab struct {
	Title  string `yaml:"title"`
	Color  string `yaml:"color,omitempty"`
	Layout Layout `yaml:"layout"`
}

// Layout defines the content of a tab (single pane or split).
type Layout struct {
	Cwd            string    `yaml:"cwd,omitempty"`
	Commands       []Command `yaml:"commands,omitempty"`
	IsFocused      bool      `yaml:"is_focused,omitempty"`
	SplitDirection string    `yaml:"split_direction,omitempty"`
	Panes          []Layout  `yaml:"panes,omitempty"`
}

// Command represents a command to execute in a pane.
type Command struct {
	Exec string `yaml:"exec"`
}

// Save writes the launch configuration to the Warp launch configurations directory.
func Save(c Config) (string, error) {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating launch config directory: %w", err)
	}

	filename := slugify(c.Name) + ".yaml"
	path := filepath.Join(dir, filename)

	data, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("marshaling launch config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("writing launch config file: %w", err)
	}

	return path, nil
}

func configDir() string {
	switch runtime.GOOS {
	case "linux":
		base := os.Getenv("XDG_DATA_HOME")
		if base == "" {
			base = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
		return filepath.Join(base, "warp-terminal", "launch_configurations")
	default: // darwin
		return filepath.Join(os.Getenv("HOME"), ".warp", "launch_configurations")
	}
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	s = strings.Trim(b.String(), "-")
	if s == "" {
		s = "launch-config"
	}
	return s
}
