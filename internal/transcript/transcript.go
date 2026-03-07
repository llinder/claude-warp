package transcript

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

// Message represents a single entry in a Claude Code transcript JSONL file.
type Message struct {
	Type    string         `json:"type"`
	Message MessageContent `json:"message"`
}

// MessageContent holds the content of a transcript message.
type MessageContent struct {
	Content json.RawMessage `json:"content"`
}

// ContentBlock represents a single content block within a message.
type ContentBlock struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Name  string `json:"name,omitempty"`
	Input struct {
		Command     string `json:"command,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"input,omitempty"`
}

// ToolResult represents a tool result content block.
type ToolResult struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
}

// BashCommand represents a bash command extracted from the transcript.
type BashCommand struct {
	Command     string
	Description string
}

// Parse reads a JSONL transcript file and returns the parsed messages.
func Parse(path string) ([]Message, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var messages []Message
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // 10MB max line
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			continue // skip malformed lines
		}
		messages = append(messages, msg)
	}
	return messages, scanner.Err()
}

// ExtractFirstPrompt returns the first user prompt text from the transcript.
func ExtractFirstPrompt(messages []Message) string {
	for _, msg := range messages {
		if msg.Type != "user" {
			continue
		}
		// Content can be a string or array of content blocks
		var text string
		if err := json.Unmarshal(msg.Message.Content, &text); err == nil {
			return text
		}
		var blocks []ContentBlock
		if err := json.Unmarshal(msg.Message.Content, &blocks); err == nil {
			for _, b := range blocks {
				if b.Type == "text" && b.Text != "" {
					return b.Text
				}
			}
		}
	}
	return ""
}

// ExtractLastResponse returns the last assistant text response from the transcript.
func ExtractLastResponse(messages []Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Type != "assistant" {
			continue
		}
		var blocks []ContentBlock
		if err := json.Unmarshal(msg.Message.Content, &blocks); err != nil {
			continue
		}
		var texts []string
		for _, b := range blocks {
			if b.Type == "text" && b.Text != "" {
				texts = append(texts, b.Text)
			}
		}
		if len(texts) > 0 {
			return strings.Join(texts, " ")
		}
	}
	return ""
}

// ExtractBashCommands returns all Bash tool invocations from the transcript.
func ExtractBashCommands(messages []Message) []BashCommand {
	var cmds []BashCommand
	for _, msg := range messages {
		if msg.Type != "assistant" {
			continue
		}
		var blocks []ContentBlock
		if err := json.Unmarshal(msg.Message.Content, &blocks); err != nil {
			continue
		}
		for _, b := range blocks {
			if b.Type == "tool_use" && b.Name == "Bash" && b.Input.Command != "" {
				cmds = append(cmds, BashCommand{
					Command:     b.Input.Command,
					Description: b.Input.Description,
				})
			}
		}
	}
	return cmds
}

// CountFilesChanged counts unique file paths from Edit/Write tool uses.
func CountFilesChanged(messages []Message) int {
	files := make(map[string]struct{})
	for _, msg := range messages {
		if msg.Type != "assistant" {
			continue
		}
		var blocks []ContentBlock
		if err := json.Unmarshal(msg.Message.Content, &blocks); err != nil {
			continue
		}
		for _, b := range blocks {
			if b.Type == "tool_use" && (b.Name == "Edit" || b.Name == "Write") {
				// Input has file_path field
				var input struct {
					FilePath string `json:"file_path"`
				}
				raw, _ := json.Marshal(b.Input)
				if json.Unmarshal(raw, &input) == nil && input.FilePath != "" {
					files[input.FilePath] = struct{}{}
				}
			}
		}
	}
	return len(files)
}
