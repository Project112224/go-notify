package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractURLs(t *testing.T) {
	svc := NewDesktopLauncherService()

	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "Visit https://github.com/june for details",
			expected: []string{"https://github.com/june"},
		},
		{
			input:    "Check www.google.com and http://example.com",
			expected: []string{"https://www.google.com", "http://example.com"},
		},
		{
			input:    "No links in this text",
			expected: nil,
		},
		{
			input:    "Duplicate https://foo.com and https://foo.com",
			expected: []string{"https://foo.com"},
		},
	}

	for _, tt := range tests {
		result := svc.ExtractURLs(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("ExtractURLs(%q) = %v; want %v", tt.input, result, tt.expected)
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("ExtractURLs(%q)[%d] = %q; want %q", tt.input, i, result[i], tt.expected[i])
			}
		}
	}
}

func TestIsSystemApp(t *testing.T) {
	if !isSystemApp("system") {
		t.Errorf("Expected system to be system app")
	}
	if !isSystemApp("notify-send") {
		t.Errorf("Expected notify-send to be system app")
	}
	if isSystemApp("firefox") {
		t.Errorf("Expected firefox not to be system app")
	}
}

func TestCleanExecLine(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "firefox %u",
			expected: "firefox",
		},
		{
			input:    "/usr/bin/code --new-window %F",
			expected: "/usr/bin/code --new-window",
		},
		{
			input:    "spotify %U",
			expected: "spotify",
		},
	}

	for _, tt := range tests {
		res := cleanExecLine(tt.input)
		if res != tt.expected {
			t.Errorf("cleanExecLine(%q) = %q; want %q", tt.input, res, tt.expected)
		}
	}
}

func TestExtractExecFromDesktopFile(t *testing.T) {
	tmpDir := t.TempDir()
	desktopFile := filepath.Join(tmpDir, "test.desktop")
	content := "[Desktop Entry]\nName=TestApp\nExec=testapp --flag %u\nType=Application\n"

	if err := os.WriteFile(desktopFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test desktop file: %v", err)
	}

	execCmd := extractExecFromDesktopFile(desktopFile)
	expected := "testapp --flag"
	if execCmd != expected {
		t.Errorf("extractExecFromDesktopFile() = %q; want %q", execCmd, expected)
	}
}
