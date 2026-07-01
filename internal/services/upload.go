package services

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/C9b3rD3vi1/SnapReport/internal/fileutil"
	"github.com/C9b3rD3vi1/SnapReport/internal/imgvalidator"
	"github.com/C9b3rD3vi1/SnapReport/internal/models"
	"github.com/C9b3rD3vi1/SnapReport/internal/service"
)

type UploadService struct {
	repo      service.UploadRepo
	uploadDir string
	maxSize   int64
}

func NewUploadService(repo service.UploadRepo, uploadDir string, maxSizeMB int64) *UploadService {
	return &UploadService{
		repo:      repo,
		uploadDir: uploadDir,
		maxSize:   maxSizeMB * 1024 * 1024,
	}
}

type uploadError struct {
	Filename string
	Err      error
}

func (s *UploadService) Create(files []*multipart.FileHeader) ([]models.Upload, []uploadError) {
	var uploads []models.Upload
	var errors []uploadError

	for _, fh := range files {
		if fh.Size > s.maxSize {
			errors = append(errors, uploadError{Filename: fh.Filename, Err: service.ErrFileTooLarge})
			continue
		}

		src, err := fh.Open()
		if err != nil {
			errors = append(errors, uploadError{Filename: fh.Filename, Err: fmt.Errorf("open: %w", err)})
			continue
		}

		header := make([]byte, 512)
		n, _ := io.ReadFull(src, header)
		header = header[:n]

		mime, err := imgvalidator.ValidateMIME(header)
		if err != nil {
			src.Close()
			errors = append(errors, uploadError{Filename: fh.Filename, Err: service.ErrInvalidFile})
			continue
		}

		id := fileutil.GenerateID()
		ext := filepath.Ext(fh.Filename)
		if e := imgvalidator.AllowedExtension(mime); e != "" {
			ext = e
		}
		filename := id + ext

		reader := io.MultiReader(bytes.NewReader(header), src)

		path, err := fileutil.Save(s.uploadDir, filename, reader)
		src.Close()
		if err != nil {
			errors = append(errors, uploadError{Filename: fh.Filename, Err: fmt.Errorf("save: %w", err)})
			continue
		}

		thumbDir := filepath.Join(s.uploadDir, "thumbnails")
		os.MkdirAll(thumbDir, 0755)
		if _, thumbErr := GenerateThumbnail(path, thumbDir); thumbErr != nil {
			slog.Warn("thumbnail generation failed", "file", filename, "error", thumbErr)
		}

		upload := models.Upload{
			ID:           id,
			Filename:     filename,
			OriginalName: fh.Filename,
			MimeType:     mime,
			Size:         fh.Size,
			Path:         path,
			CreatedAt:    time.Now(),
		}
		upload.SetThumbnailURL()

		if err := s.repo.Insert(&upload); err != nil {
			fileutil.Remove(path)
			errors = append(errors, uploadError{Filename: fh.Filename, Err: fmt.Errorf("db: %w", err)})
			continue
		}

		uploads = append(uploads, upload)
		slog.Info("file uploaded", "id", upload.ID, "filename", upload.OriginalName, "size", upload.Size)
	}

	return uploads, errors
}

func (s *UploadService) List() ([]models.Upload, error) {
	uploads, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("list uploads: %w", err)
	}
	if uploads == nil {
		uploads = []models.Upload{}
	}
	for i := range uploads {
		uploads[i].SetThumbnailURL()
	}
	return uploads, nil
}

func (s *UploadService) Delete(id string) error {
	upload, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("find upload: %w", err)
	}
	if upload == nil {
		return service.ErrNotFound
	}

	if err := fileutil.Remove(upload.Path); err != nil {
		slog.Warn("failed to remove file", "path", upload.Path, "error", err)
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete upload record: %w", err)
	}

	slog.Info("file deleted", "id", id)
	return nil
}
