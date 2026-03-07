package notify

import "testing"

func TestSanitize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"clean text", "Hello World", "Hello World"},
		{"escape sequence", "Hello\033[31mWorld", "Hello[31mWorld"},
		{"null bytes", "Hello\x00World", "HelloWorld"},
		{"newlines", "Hello\nWorld", "HelloWorld"},
		{"tabs", "Hello\tWorld", "HelloWorld"},
		{"bell", "Hello\x07World", "HelloWorld"},
		{"DEL", "Hello\x7FWorld", "HelloWorld"},
		{"unicode preserved", "Hello World", "Hello World"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitize(tt.input)
			if got != tt.want {
				t.Errorf("sanitize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is too long", 10, "this is..."},
		{"ab", 2, "ab"},
		{"abc", 2, "ab"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}
