package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveAndList(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	w := Workflow{
		Name:        "Test Workflow",
		Command:     "echo {{message}}",
		Description: "A test workflow",
		Tags:        []string{"test"},
		Arguments: []Argument{
			{Name: "message", Description: "The message to print", DefaultValue: "hello"},
		},
	}

	path, err := Save(w, false)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if filepath.Ext(path) != ".yaml" {
		t.Errorf("expected .yaml extension, got %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading saved file: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("saved file is empty")
	}

	workflows, err := List(false)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(workflows) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(workflows))
	}
	if workflows[0].Name != "Test Workflow" {
		t.Errorf("expected name 'Test Workflow', got %q", workflows[0].Name)
	}
	if workflows[0].Command != "echo {{message}}" {
		t.Errorf("expected command 'echo {{message}}', got %q", workflows[0].Command)
	}
}

func TestSaveRepoScoped(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a fake git repo
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0755)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	w := Workflow{
		Name:    "Repo Workflow",
		Command: "make test",
	}

	path, err := Save(w, true)
	if err != nil {
		t.Fatalf("Save(repo) error: %v", err)
	}

	// Resolve symlinks (macOS /var -> /private/var)
	realTmp, _ := filepath.EvalSymlinks(tmpDir)
	wantPrefix := filepath.Join(realTmp, ".warp", "workflows")
	if !strings.HasPrefix(path, wantPrefix) {
		t.Errorf("expected repo-scoped path under %s, got %s", wantPrefix, path)
	}

	workflows, err := List(true)
	if err != nil {
		t.Fatalf("List(repo) error: %v", err)
	}
	if len(workflows) != 1 {
		t.Fatalf("expected 1 repo workflow, got %d", len(workflows))
	}
}

func TestSaveRepoScoped_NotInGit(t *testing.T) {
	origDir, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(origDir)

	_, err := Save(Workflow{Name: "test", Command: "echo"}, true)
	if err == nil {
		t.Error("expected error when not in git repo")
	}
}

func TestListEmpty(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	workflows, err := List(false)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(workflows) != 0 {
		t.Errorf("expected 0 workflows, got %d", len(workflows))
	}
}

func TestListSkipsBadFiles(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	wfDir := filepath.Join(tmpDir, ".warp", "workflows")
	os.MkdirAll(wfDir, 0755)

	// Write a valid workflow
	os.WriteFile(filepath.Join(wfDir, "good.yaml"), []byte("name: Good\ncommand: echo\n"), 0644)
	// Write an invalid YAML file
	os.WriteFile(filepath.Join(wfDir, "bad.yaml"), []byte("{{invalid yaml"), 0644)
	// Write a non-yaml file (should be skipped)
	os.WriteFile(filepath.Join(wfDir, "readme.txt"), []byte("not a workflow"), 0644)

	workflows, err := List(false)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(workflows) != 1 {
		t.Errorf("expected 1 valid workflow, got %d", len(workflows))
	}
}

func TestPersonalDir(t *testing.T) {
	t.Setenv("HOME", "/fakehome")
	dir := PersonalDir()
	if dir != "/fakehome/.warp/workflows" {
		t.Errorf("PersonalDir() = %q, want /fakehome/.warp/workflows", dir)
	}
}

func TestRepoDir_NotInGit(t *testing.T) {
	origDir, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(origDir)

	dir := RepoDir()
	if dir != "" {
		t.Errorf("RepoDir() = %q, want empty string", dir)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"git: push & pull", "git-push-pull"},
		{"", "workflow"},
		{"---", "workflow"},
		{"Simple", "simple"},
		{"Multiple   Spaces", "multiple-spaces"},
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
