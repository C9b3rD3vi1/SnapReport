package services

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/anomalyco/SnapReport/internal/models"
	"github.com/anomalyco/SnapReport/internal/repository"
	imgvalidator "github.com/anomalyco/SnapReport/internal/services/image"
	"github.com/anomalyco/SnapReport/internal/utils"
)

type UploadService struct {
	repo      *repository.UploadRepository
	uploadDir string
	maxSize   int64
}

func NewUploadService(repo *repository.UploadRepository, uploadDir string, maxSizeMB int64) *UploadService {
	return &UploadService{
		repo:      repo,
		uploadDir: uploadDir,
		maxSize:   maxSizeMB * 1024 * 1024,
	}
}

func (s *UploadService) Create(files []*multipart.FileHeader) ([]models.Upload, error) {
	var uploads []models.Upload

	for _, fh := range files {
		if fh.Size > s.maxSize {
			return nil, fmt.Errorf("file %s exceeds maximum size of %d MB", fh.Filename, s.maxSize/(1024*1024))
		}

		src, err := fh.Open()
		if err != nil {
			return nil, fmt.Errorf("open file %s: %w", fh.Filename, err)
		}

		header := make([]byte, 512)
		n, err := io.ReadFull(src, header)
		if err != nil && err != io.ErrUnexpectedEOF {
			src.Close()
			return nil, fmt.Errorf("read file header: %w", err)
		}
		header = header[:n]

		mime, err := imgvalidator.ValidateMIME(header)
		if err != nil {
			src.Close()
			return nil, fmt.Errorf("file %s: %w", fh.Filename, err)
		}

		id := utils.GenerateID()
		ext := filepath.Ext(fh.Filename)
		if e := imgvalidator.AllowedExtension(mime); e != "" {
			ext = e
		}
		filename := id + ext

		reader := io.MultiReader(bytes.NewReader(header), src)

		path, err := utils.SaveFile(s.uploadDir, filename, reader)
		src.Close()
		if err != nil {
			return nil, fmt.Errorf("save file %s: %w", filename, err)
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

		if err := s.repo.Insert(&upload); err != nil {
			utils.RemoveFile(path)
			return nil, fmt.Errorf("save upload record: %w", err)
		}

		uploads = append(uploads, upload)
		slog.Info("file uploaded", "id", upload.ID, "filename", upload.OriginalName, "size", upload.Size)
	}

	return uploads, nil
}

func (s *UploadService) List() ([]models.Upload, error) {
	uploads, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("list uploads: %w", err)
	}
	if uploads == nil {
		uploads = []models.Upload{}
	}
	return uploads, nil
}

func (s *UploadService) Delete(id string) error {
	upload, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("find upload: %w", err)
	}
	if upload == nil {
		return fmt.Errorf("upload not found")
	}

	if err := utils.RemoveFile(upload.Path); err != nil {
		slog.Warn("failed to remove file", "path", upload.Path, "error", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete upload record: %w", err)
	}

	slog.Info("file deleted", "id", id)
	return nil
}
