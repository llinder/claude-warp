package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Workflow represents a Warp workflow YAML file.
type Workflow struct {
	Name        string     `yaml:"name"`
	Command     string     `yaml:"command"`
	Description string     `yaml:"description,omitempty"`
	Tags        []string   `yaml:"tags,omitempty"`
	Shells      []string   `yaml:"shells,omitempty"`
	Arguments   []Argument `yaml:"arguments,omitempty"`
}

// Argument represents a workflow argument.
type Argument struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description,omitempty"`
	DefaultValue any    `yaml:"default_value"`
}

// Save writes the workflow to the appropriate directory.
// If repoScoped is true, saves to .warp/workflows/ in the git repo root.
// Otherwise saves to the user's personal Warp workflows directory.
func Save(w Workflow, repoScoped bool) (string, error) {
	dir, err := targetDir(repoScoped)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating workflow directory: %w", err)
	}

	filename := slugify(w.Name) + ".yaml"
	path := filepath.Join(dir, filename)

	data, err := yaml.Marshal(w)
	if err != nil {
		return "", fmt.Errorf("marshaling workflow: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("writing workflow file: %w", err)
	}

	return path, nil
}

// List returns all workflow files in the given directory.
func List(repoScoped bool) ([]Workflow, error) {
	dir, err := targetDir(repoScoped)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var workflows []Workflow
	for _, e := range entries {
		if e.IsDir() || (!strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml")) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var w Workflow
		if err := yaml.Unmarshal(data, &w); err != nil {
			continue
		}
		workflows = append(workflows, w)
	}
	return workflows, nil
}

// PersonalDir returns the user's personal Warp workflows directory.
func PersonalDir() string {
	switch runtime.GOOS {
	case "linux":
		base := os.Getenv("XDG_DATA_HOME")
		if base == "" {
			base = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
		return filepath.Join(base, "warp-terminal", "workflows")
	default: // darwin
		return filepath.Join(os.Getenv("HOME"), ".warp", "workflows")
	}
}

// RepoDir returns the repo-scoped Warp workflows directory, or empty string if not in a git repo.
func RepoDir() string {
	root := gitRoot()
	if root == "" {
		return ""
	}
	return filepath.Join(root, ".warp", "workflows")
}

func targetDir(repoScoped bool) (string, error) {
	if repoScoped {
		dir := RepoDir()
		if dir == "" {
			return "", fmt.Errorf("not in a git repository")
		}
		return dir, nil
	}
	return PersonalDir(), nil
}

func gitRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "workflow"
	}
	return s
}
