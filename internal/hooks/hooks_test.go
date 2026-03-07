package hooks

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSessionStart_InWarp(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "WarpTerminal")
	// Use temp HOME so workflow discovery doesn't find real files
	t.Setenv("HOME", t.TempDir())

	output := captureStdout(t, func() {
		if err := SessionStart(); err != nil {
			t.Fatalf("SessionStart() error: %v", err)
		}
	})

	if !strings.Contains(output, "systemMessage") {
		t.Error("expected systemMessage in output")
	}
	if !strings.Contains(output, "Warp terminal detected") {
		t.Error("expected Warp detection message")
	}
	if !strings.Contains(output, "save-workflow") {
		t.Error("expected workflow instructions")
	}
	if !strings.Contains(output, "save-launch") {
		t.Error("expected launch config instructions")
	}
}

func TestSessionStart_NotInWarp(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "iTerm2")

	output := captureStdout(t, func() {
		if err := SessionStart(); err != nil {
			t.Fatalf("SessionStart() error: %v", err)
		}
	})

	if !strings.Contains(output, "not running in Warp") {
		t.Error("expected not-in-Warp message")
	}
}

func TestSessionStart_WithExistingWorkflows(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "WarpTerminal")
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Create a fake workflow
	wfDir := filepath.Join(tmpDir, ".warp", "workflows")
	os.MkdirAll(wfDir, 0755)
	os.WriteFile(filepath.Join(wfDir, "test.yaml"), []byte("name: Test WF\ncommand: echo hi\ndescription: A test\n"), 0644)

	output := captureStdout(t, func() {
		if err := SessionStart(); err != nil {
			t.Fatalf("SessionStart() error: %v", err)
		}
	})

	if !strings.Contains(output, "Personal workflows") {
		t.Error("expected personal workflows listing")
	}
}

func TestStop(t *testing.T) {
	// Create a sample transcript
	transcriptDir := t.TempDir()
	transcriptPath := filepath.Join(transcriptDir, "transcript.jsonl")
	transcript := `{"type":"user","message":{"content":"Fix the bug"}}
{"type":"assistant","message":{"content":[{"type":"text","text":"Done fixing the bug."}]}}
`
	os.WriteFile(transcriptPath, []byte(transcript), 0644)

	// Feed the transcript path via stdin
	input := `{"transcript_path":"` + transcriptPath + `"}`
	setStdin(t, input)

	// Stop will try to send notification to /dev/tty which may fail in test,
	// but the function itself should not error
	err := Stop()
	if err != nil {
		t.Fatalf("Stop() error: %v", err)
	}
}

func TestStop_NoTranscript(t *testing.T) {
	setStdin(t, `{}`)
	err := Stop()
	if err != nil {
		t.Fatalf("Stop() error: %v", err)
	}
}

func TestNotification(t *testing.T) {
	setStdin(t, `{"message":"Build complete"}`)
	err := Notification()
	if err != nil {
		t.Fatalf("Notification() error: %v", err)
	}
}

func TestNotification_EmptyMessage(t *testing.T) {
	setStdin(t, `{}`)
	err := Notification()
	if err != nil {
		t.Fatalf("Notification() error: %v", err)
	}
}

func TestBuildStopMessage(t *testing.T) {
	transcriptDir := t.TempDir()

	tests := []struct {
		name       string
		transcript string
		wantSubstr string
	}{
		{
			name:       "empty path",
			transcript: "",
			wantSubstr: "Task completed",
		},
		{
			name: "with prompt and response",
			transcript: `{"type":"user","message":{"content":"Deploy the app"}}
{"type":"assistant","message":{"content":[{"type":"text","text":"Deployed successfully."}]}}
`,
			wantSubstr: "Deploy the app",
		},
		{
			name: "with bash commands",
			transcript: `{"type":"user","message":{"content":"Run tests"}}
{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","input":{"command":"go test ./...","description":"run tests"}},{"type":"text","text":"All tests pass."}]}}
`,
			wantSubstr: "1 commands run",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string
			if tt.transcript != "" {
				path = filepath.Join(transcriptDir, tt.name+".jsonl")
				os.WriteFile(path, []byte(tt.transcript), 0644)
			}

			msg := buildStopMessage(path)
			if !strings.Contains(msg, tt.wantSubstr) {
				t.Errorf("buildStopMessage() = %q, want substring %q", msg, tt.wantSubstr)
			}
		})
	}
}

func TestSaveWorkflow(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	err := SaveWorkflow(SaveWorkflowOpts{
		Name:        "Test Deploy",
		Command:     "kubectl apply -f {{file}}",
		Description: "Apply a k8s manifest",
		ArgSpecs:    []string{"file:Manifest file path:manifest.yaml"},
		Tags:        []string{"k8s"},
	})
	if err != nil {
		t.Fatalf("SaveWorkflow() error: %v", err)
	}
}

func TestSaveWorkflow_RepoNotInGit(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	origDir, _ := os.Getwd()
	os.Chdir(t.TempDir())
	defer os.Chdir(origDir)

	err := SaveWorkflow(SaveWorkflowOpts{
		Name:       "test",
		Command:    "echo hi",
		RepoScoped: true,
	})
	if err == nil {
		t.Error("expected error for --repo outside git repo")
	}
}

func TestSaveLaunch(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	err := SaveLaunch(SaveLaunchOpts{
		Name:     "Dev Env",
		TabSpecs: []string{"Server:/tmp/project:npm run dev", "Tests:/tmp/project:npm test"},
	})
	if err != nil {
		t.Fatalf("SaveLaunch() error: %v", err)
	}
}

func TestParseArgSpec(t *testing.T) {
	tests := []struct {
		spec        string
		wantName    string
		wantDesc    string
		wantDefault any
		wantErr     bool
	}{
		{"name:desc:default", "name", "desc", "default", false},
		{"name:desc:~", "name", "desc", nil, false},
		{"name:desc", "name", "desc", nil, false},
		{"name", "name", "", nil, false},
		{"", "", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.spec, func(t *testing.T) {
			arg, err := parseArgSpec(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseArgSpec(%q) error = %v, wantErr %v", tt.spec, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if arg.Name != tt.wantName {
				t.Errorf("name = %q, want %q", arg.Name, tt.wantName)
			}
			if arg.Description != tt.wantDesc {
				t.Errorf("description = %q, want %q", arg.Description, tt.wantDesc)
			}
			if arg.DefaultValue != tt.wantDefault {
				t.Errorf("default = %v, want %v", arg.DefaultValue, tt.wantDefault)
			}
		})
	}
}

func TestParseTabSpec(t *testing.T) {
	tests := []struct {
		spec      string
		wantTitle string
		wantCwd   string
		wantCmd   string
		wantErr   bool
	}{
		{"Server:/tmp:npm start", "Server", "/tmp", "npm start", false},
		{"Server:/tmp", "Server", "/tmp", "", false},
		{"Server", "Server", "", "", false},
		{"", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.spec, func(t *testing.T) {
			tab, err := parseTabSpec(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseTabSpec(%q) error = %v, wantErr %v", tt.spec, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if tab.Title != tt.wantTitle {
				t.Errorf("title = %q, want %q", tab.Title, tt.wantTitle)
			}
			if tab.Layout.Cwd != tt.wantCwd {
				t.Errorf("cwd = %q, want %q", tab.Layout.Cwd, tt.wantCwd)
			}
			cmd := ""
			if len(tab.Layout.Commands) > 0 {
				cmd = tab.Layout.Commands[0].Exec
			}
			if cmd != tt.wantCmd {
				t.Errorf("command = %q, want %q", cmd, tt.wantCmd)
			}
		})
	}
}

// captureStdout redirects stdout to capture output from a function.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// setStdin replaces stdin with a string reader for the duration of the test.
func setStdin(t *testing.T, s string) {
	t.Helper()
	old := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte(s))
	w.Close()
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = old })
}
