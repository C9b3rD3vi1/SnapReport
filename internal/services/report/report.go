package report

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/C9b3rD3vi1/SnapReport/internal/fileutil"
	"github.com/C9b3rD3vi1/SnapReport/internal/models"
	"github.com/C9b3rD3vi1/SnapReport/internal/pdf"
	"github.com/C9b3rD3vi1/SnapReport/internal/repository"
	"github.com/C9b3rD3vi1/SnapReport/internal/service"
)

type Service struct {
	reportRepo service.ReportRepo
	uploadRepo service.UploadRepo
	blockRepo  *repository.BlockRepository
	pdfGen     *pdf.Generator
	pdfDir     string
}

func NewService(reportRepo service.ReportRepo, uploadRepo service.UploadRepo, blockRepo *repository.BlockRepository, pdfGen *pdf.Generator, pdfDir string) *Service {
	return &Service{
		reportRepo: reportRepo,
		uploadRepo: uploadRepo,
		blockRepo:  blockRepo,
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
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Notes          string `json:"notes"`
	OrderIndex     int    `json:"order_index"`
	Category       string `json:"category"`
	Priority       string `json:"priority"`
	Severity       string `json:"severity"`
	Status         string `json:"status"`
	Recommendation string `json:"recommendation"`
}

func (s *Service) Create(req CreateRequest) (*models.Report, error) {
	if req.Title == "" {
		return nil, service.ErrTitleMissing
	}
	if len(req.Uploads) == 0 {
		return nil, service.ErrNoUploads
	}

	id := fileutil.GenerateID()
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
		if err := s.uploadRepo.UpdateMetadata(u.ID, u.Title, u.Description, u.Notes,
			u.Category, u.Priority, u.Severity, u.Status, u.Recommendation,
			u.OrderIndex, id); err != nil {
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

func (s *Service) PreviewHTML(req CreateRequest) (string, error) {
	uploads, err := s.uploadRepo.FindAll()
	if err != nil {
		return "", fmt.Errorf("fetch uploads: %w", err)
	}

	uploadMap := map[string]models.Upload{}
	for _, u := range uploads {
		uploadMap[u.ID] = u
	}

	var screenshots []pdf.ScreenshotData
	for i, u := range req.Uploads {
		up, ok := uploadMap[u.ID]
		if !ok {
			continue
		}
		title := u.Title
		if title == "" {
			title = fmt.Sprintf("Finding %d", i+1)
		}
		screenshots = append(screenshots, pdf.ScreenshotData{
			ImagePath:      "/uploads/" + up.Filename,
			Title:          title,
			Description:    u.Description,
			FigureLabel:    fmt.Sprintf("Figure %d", i+1),
			Category:       u.Category,
			Priority:       u.Priority,
			Severity:       u.Severity,
			Status:         u.Status,
			Recommendation: u.Recommendation,
		})
	}

	data := pdf.ReportData{
		Title:          req.Title,
		Project:        req.Project,
		Company:        req.Company,
		Author:         req.Author,
		Version:        req.Version,
		Date:           time.Now().Format("January 2, 2006"),
		Classification: "Internal",
		ReportID:       "PREVIEW",
		Status:         "draft",
		Watermark:      "PREVIEW",
		Screenshots:    screenshots,
		Summary: pdf.SummaryData{
			TotalFindings:      len(screenshots),
			TotalImages:        len(screenshots),
			HasRecommendations: false,
			ReadingTime:        fmt.Sprintf("%d min", max(1, (len(screenshots)+2)/2)),
		},
	}

	for _, u := range req.Uploads {
		switch u.Priority {
		case "High", "Critical":
			data.Summary.HighCount++
		case "Medium":
			data.Summary.MediumCount++
		case "Low":
			data.Summary.LowCount++
		}
	}

	if s.pdfGen != nil {
		return s.pdfGen.RenderHTML(data)
	}
	return "", fmt.Errorf("pdf generator not available")
}

func (s *Service) generatePDF(report *models.Report, uploads []models.Upload) error {
	high, med, low := 0, 0, 0
	catSet := map[string]bool{}
	hasRecs := false

	for _, u := range uploads {
		switch u.Priority {
		case "High", "Critical":
			high++
		case "Medium":
			med++
		case "Low":
			low++
		}
		if u.Category != "" {
			catSet[u.Category] = true
		}
		if u.Recommendation != "" {
			hasRecs = true
		}
	}

	var cats []string
	for c := range catSet {
		cats = append(cats, c)
	}
	if cats == nil {
		cats = []string{}
	}

	totalFindings := len(uploads)
	readingTime := fmt.Sprintf("%d min", max(1, (totalFindings+2)/2))

	blocks, _ := s.blockRepo.FindByReportID(report.ID)
	umap := map[string]models.Upload{}
	for _, u := range uploads {
		umap[u.ID] = u
	}
	var pdfBlocks []pdf.BlockData
	for _, b := range blocks {
		pdfBlocks = append(pdfBlocks, blockToPDF(b, umap))
	}

	data := pdf.ReportData{
		Title:          report.Title,
		Project:        report.Project,
		Company:        report.Company,
		Author:         report.Author,
		Version:        report.Version,
		Date:           report.CreatedAt.Format("January 2, 2006"),
		Classification: "Internal",
		ReportID:       report.ID[:8],
		Status:         report.Status,
		Watermark:      "",
		Blocks:         pdfBlocks,
		Summary: pdf.SummaryData{
			TotalFindings:      totalFindings,
			TotalImages:        totalFindings,
			HighCount:          high,
			MediumCount:        med,
			LowCount:           low,
			Categories:         cats,
			ReadingTime:        readingTime,
			HasRecommendations: hasRecs,
		},
	}

	for i, u := range uploads {
		caption := fmt.Sprintf("Figure %d", i+1)
		if u.Title != "" {
			caption = fmt.Sprintf("Figure %d: %s", i+1, u.Title)
		}
		data.Screenshots = append(data.Screenshots, pdf.ScreenshotData{
			ImagePath:      u.Path,
			Title:          u.Title,
			Description:    u.Description,
			FigureLabel:    caption,
			Category:       u.Category,
			Priority:       u.Priority,
			Severity:       u.Severity,
			Status:         u.Status,
			Recommendation: u.Recommendation,
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

func (s *Service) List() ([]models.Report, error) {
	reports, err := s.reportRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	if reports == nil {
		reports = []models.Report{}
	}
	return reports, nil
}

func (s *Service) Delete(id string) error {
	report, err := s.reportRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("find report: %w", err)
	}
	if report == nil {
		return service.ErrNotFound
	}
	if report.PDFPath != "" {
		fileutil.Remove(report.PDFPath)
	}
	return s.reportRepo.Delete(id)
}

func blockToPDF(b models.Block, uploadMap map[string]models.Upload) pdf.BlockData {
	pb := pdf.BlockData{Type: string(b.Type)}

	switch b.Type {
	case models.BlockFinding:
		var c struct {
			ScreenshotID   string `json:"screenshot_id"`
			Title          string `json:"title"`
			Description    string `json:"description"`
			Category       string `json:"category"`
			Priority       string `json:"priority"`
			Severity       string `json:"severity"`
			Status         string `json:"status"`
			Recommendation string `json:"recommendation"`
		}
		json.Unmarshal([]byte(b.Content), &c)
		pb.FindingTitle = c.Title
		pb.FindingDesc = c.Description
		pb.Category = c.Category
		pb.Priority = c.Priority
		pb.Severity = c.Severity
		pb.Status = c.Status
		pb.Recommendation = c.Recommendation
		if up, ok := uploadMap[c.ScreenshotID]; ok {
			pb.FindingImage = up.Path
		}

	case models.BlockNote:
		var c struct{ HTML string `json:"html"` }
		json.Unmarshal([]byte(b.Content), &c)
		pb.NoteContent = c.HTML

	case models.BlockWarning, models.BlockTip, models.BlockImportant:
		var c struct{ Text string `json:"text"` }
		json.Unmarshal([]byte(b.Content), &c)
		pb.MessageText = c.Text

	case models.BlockDivider:
		var c struct{ Title string `json:"title"` }
		json.Unmarshal([]byte(b.Content), &c)
		pb.DividerTitle = c.Title

	case models.BlockChecklist:
		var c struct {
			Items []struct {
				Text    string `json:"text"`
				Checked bool   `json:"checked"`
			} `json:"items"`
		}
		json.Unmarshal([]byte(b.Content), &c)
		for _, item := range c.Items {
			pb.ChecklistItems = append(pb.ChecklistItems, pdf.ChecklistItemData{Text: item.Text, Checked: item.Checked})
		}
	}

	return pb
}

func (s *Service) Get(id string) (*models.Report, []models.Upload, error) {
	report, err := s.reportRepo.FindByID(id)
	if err != nil {
		return nil, nil, fmt.Errorf("find report: %w", err)
	}
	if report == nil {
		return nil, nil, service.ErrNotFound
	}

	uploads, err := s.uploadRepo.FindByReportID(id)
	if err != nil {
		return nil, nil, fmt.Errorf("find report uploads: %w", err)
	}
	for i := range uploads {
		uploads[i].SetThumbnailURL()
	}
	return report, uploads, nil
}
