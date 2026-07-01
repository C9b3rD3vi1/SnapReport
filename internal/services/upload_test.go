package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anomalyco/SnapReport/internal/repository"
)

func setupTest(t *testing.T) (*UploadService, string) {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	sqlite, err := repository.NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("NewSQLite() error = %v", err)
	}
	t.Cleanup(func() { sqlite.Close() })

	if err := sqlite.Migrate(); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	uploadDir := filepath.Join(dir, "uploads")
	os.MkdirAll(uploadDir, 0755)

	repo := repository.NewUploadRepository(sqlite.DB())
	svc := NewUploadService(repo, uploadDir, 10)

	return svc, dir
}

func TestUploadService_List(t *testing.T) {
	svc, _ := setupTest(t)

	uploads, err := svc.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(uploads) != 0 {
		t.Errorf("expected empty list, got %d items", len(uploads))
	}
}

func TestUploadService_DeleteNonexistent(t *testing.T) {
	svc, _ := setupTest(t)

	err := svc.Delete("nonexistent-id")
	if err == nil {
		t.Fatal("expected error for nonexistent upload")
	}
}

func TestUploadService_ListAfterDelete(t *testing.T) {
	svc, _ := setupTest(t)

	uploads, err := svc.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(uploads) != 0 {
		t.Errorf("expected empty list, got %d items", len(uploads))
	}
}
