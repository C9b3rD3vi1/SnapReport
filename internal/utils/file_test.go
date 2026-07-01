package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()

	if id1 == "" {
		t.Error("expected non-empty id")
	}
	if id1 == id2 {
		t.Error("expected unique ids")
	}
	if len(id1) != 32 {
		t.Errorf("expected 32 hex chars, got %d", len(id1))
	}
}

func TestSaveFile(t *testing.T) {
	dir := t.TempDir()
	content := strings.NewReader("hello world")

	path, err := SaveFile(dir, "test.txt", content)
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}

	if filepath.Base(path) != "test.txt" {
		t.Errorf("expected filename test.txt, got %s", filepath.Base(path))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got %q", string(data))
	}
}

func TestRemoveFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("data"), 0644)

	if err := RemoveFile(path); err != nil {
		t.Fatalf("RemoveFile() error = %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestRemoveFileNonExistent(t *testing.T) {
	err := RemoveFile("/nonexistent/path/file.txt")
	if err != nil {
		t.Errorf("expected no error for non-existent file, got %v", err)
	}
}
