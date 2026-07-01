package report

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/anomalyco/SnapReport/internal/models"
	"github.com/anomalyco/SnapReport/internal/pdf"
	"github.com/anomalyco/SnapReport/internal/repository"
	"github.com/anomalyco/SnapReport/internal/utils"
)

type Service struct {
	reportRepo *repository.ReportRepository
	uploadRepo *repository.UploadRepository
	pdfGen     *pdf.Generator
	pdfDir     string
}

func NewService(reportRepo *repository.ReportRepository, uploadRepo *repository.UploadRepository, pdfGen *pdf.Generator, pdfDir string) *Service {
	return &Service{
		reportRepo: reportRepo,
		uploadRepo: uploadRepo,
		pdfGen:     pdfGen,
		pdfDir:     pdfDir,
	}
}

type CreateRequest struct {
	Title   string                `json:"title"`
	Project string                `json:"project"`
	Company string                `json:"company"`
	Author  string                `json:"author"`
	Version string                `json:"version"`
	Uploads []UploadMetadataInput `json:"uploads"`
}

type UploadMetadataInput struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
	OrderIndex  int    `json:"order_index"`
}

func (s *Service) Create(req CreateRequest) (*models.Report, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("report title is required")
	}
	if len(req.Uploads) == 0 {
		return nil, fmt.Errorf("at least one upload is required")
	}

	id := utils.GenerateID()
	now := time.Now()

	report := &models.Report{
		ID:        id,
		Title:     req.Title,
		Project:   req.Project,
		Company:   req.Company,
		Author:    req.Author,
		Version:   req.Version,
		Status:    "draft",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.reportRepo.Insert(report); err != nil {
		return nil, fmt.Errorf("create report: %w", err)
	}

	for _, u := range req.Uploads {
		if err := s.uploadRepo.UpdateMetadata(u.ID, u.Title, u.Description, u.Notes, u.OrderIndex, id); err != nil {
			return nil, fmt.Errorf("update upload metadata: %w", err)
		}
	}

	uploads, err := s.uploadRepo.FindByReportID(id)
	if err != nil {
		return nil, fmt.Errorf("fetch uploads for pdf: %w", err)
	}

	if s.pdfGen != nil {
		if err := s.generatePDF(report, uploads); err != nil {
			return nil, fmt.Errorf("generate pdf: %w", err)
		}
	}

	slog.Info("report created",
		"id", id,
		"title", req.Title,
		"uploads", len(req.Uploads),
		"status", report.Status,
	)
	return report, nil
}

func (s *Service) generatePDF(report *models.Report, uploads []models.Upload) error {
	data := pdf.ReportData{
		Title:   report.Title,
		Project: report.Project,
		Company: report.Company,
		Author:  report.Author,
		Version: report.Version,
		Date:    report.CreatedAt.Format("January 2, 2006"),
	}

	for _, u := range uploads {
		data.Screenshots = append(data.Screenshots, pdf.ScreenshotData{
			ImagePath:   u.Path,
			Title:       u.Title,
			Description: u.Description,
		})
	}

	pdfPath := filepath.Join(s.pdfDir, report.ID+".pdf")
	if err := s.pdfGen.Generate(data, pdfPath); err != nil {
		return fmt.Errorf("generate pdf file: %w", err)
	}

	report.PDFPath = pdfPath
	report.Status = "completed"

	if err := s.reportRepo.UpdatePDFPath(report.ID, pdfPath); err != nil {
		return fmt.Errorf("update report pdf path: %w", err)
	}

	slog.Info("pdf generated", "report_id", report.ID, "path", pdfPath)
	return nil
}

func (s *Service) Get(id string) (*models.Report, []models.Upload, error) {
	report, err := s.reportRepo.FindByID(id)
	if err != nil {
		return nil, nil, fmt.Errorf("find report: %w", err)
	}
	if report == nil {
		return nil, nil, fmt.Errorf("report not found")
	}

	uploads, err := s.uploadRepo.FindByReportID(id)
	if err != nil {
		return nil, nil, fmt.Errorf("find report uploads: %w", err)
	}

	return report, uploads, nil
}
