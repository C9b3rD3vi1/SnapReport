package fileutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()
	if id1 == "" || id1 == id2 {
		t.Error("expected unique non-empty ids")
	}
	if len(id1) != 32 {
		t.Errorf("expected 32 hex chars, got %d", len(id1))
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	content := strings.NewReader("hello world")
	path, err := Save(dir, "test.txt", content)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if filepath.Base(path) != "test.txt" {
		t.Errorf("expected test.txt, got %s", filepath.Base(path))
	}
	data, _ := os.ReadFile(path)
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got %q", string(data))
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("data"), 0644)
	if err := Remove(path); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestRemoveNonExistent(t *testing.T) {
	if err := Remove("/nonexistent/path/file.txt"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
