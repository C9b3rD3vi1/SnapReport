package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/C9b3rD3vi1/SnapReport/internal/repository"
)

func setupTest(t *testing.T) *Service {
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

	reportRepo := repository.NewReportRepository(sqlite.DB())
	uploadRepo := repository.NewUploadRepository(sqlite.DB())

	return NewService(reportRepo, uploadRepo, nil, "")
}

func TestService_CreateMissingTitle(t *testing.T) {
	svc := setupTest(t)

	_, err := svc.Create(CreateRequest{
		Title:   "",
		Uploads: []UploadMetadataInput{{ID: "test-id"}},
	})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestService_CreateMissingUploads(t *testing.T) {
	svc := setupTest(t)

	_, err := svc.Create(CreateRequest{
		Title:   "Test Report",
		Uploads: []UploadMetadataInput{},
	})
	if err == nil {
		t.Fatal("expected error for empty uploads")
	}
}

func TestService_CreateWithValidData(t *testing.T) {
	svc := setupTest(t)

	req := CreateRequest{
		Title:   "Test Report",
		Project: "Test Project",
		Company: "Test Corp",
		Author:  "Tester",
		Version: "1.0",
		Uploads: []UploadMetadataInput{
			{
				ID:          "nonexistent-upload-id",
				Title:       "Screenshot 1",
				Description: "First screenshot",
				Notes:       "",
				OrderIndex:  0,
			},
		},
	}

	_, err := svc.Create(req)
	if err == nil {
		t.Fatal("expected error for nonexistent upload")
	}
}

func TestService_GetNonexistent(t *testing.T) {
	svc := setupTest(t)

	_, _, err := svc.Get("nonexistent-id")
	if err == nil {
		t.Fatal("expected error for nonexistent report")
	}
}

func TestService_FindByIDReturnsNilForMissing(t *testing.T) {
	svc := setupTest(t)

	report, err := svc.reportRepo.FindByID("nonexistent")
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if report != nil {
		t.Fatal("expected nil for nonexistent report")
	}

	uploads, err := svc.uploadRepo.FindByReportID("nonexistent")
	if err != nil {
		t.Fatalf("FindByReportID() error = %v", err)
	}
	if len(uploads) != 0 {
		t.Fatalf("expected 0 uploads, got %d", len(uploads))
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
