package launch

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	config := Config{
		Name:              "Dev Environment",
		ActiveWindowIndex: 0,
		Windows: []Window{
			{
				ActiveTabIndex: 0,
				Tabs: []Tab{
					{
						Title: "Server",
						Color: "Green",
						Layout: Layout{
							Cwd:       "/tmp/project",
							Commands:  []Command{{Exec: "npm run dev"}},
							IsFocused: true,
						},
					},
					{
						Title: "Tests",
						Color: "Yellow",
						Layout: Layout{
							Cwd:      "/tmp/project",
							Commands: []Command{{Exec: "npm test -- --watch"}},
						},
					},
				},
			},
		},
	}

	path, err := Save(config)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading saved file: %v", err)
	}

	var loaded Config
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("unmarshaling saved file: %v", err)
	}

	if loaded.Name != "Dev Environment" {
		t.Errorf("expected name 'Dev Environment', got %q", loaded.Name)
	}
	if len(loaded.Windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(loaded.Windows))
	}
	if len(loaded.Windows[0].Tabs) != 2 {
		t.Fatalf("expected 2 tabs, got %d", len(loaded.Windows[0].Tabs))
	}
	if loaded.Windows[0].Tabs[0].Title != "Server" {
		t.Errorf("expected tab title 'Server', got %q", loaded.Windows[0].Tabs[0].Title)
	}
	if loaded.Windows[0].Tabs[0].Layout.Cwd != "/tmp/project" {
		t.Errorf("expected cwd '/tmp/project', got %q", loaded.Windows[0].Tabs[0].Layout.Cwd)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"My Dev Env", "my-dev-env"},
		{"", "launch-config"},
		{"test", "test"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := slugify(tt.input)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
