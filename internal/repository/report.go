package repository

import (
	"database/sql"
	"fmt"
	"github.com/anomalyco/SnapReport/internal/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Insert(report *models.Report) error {
	query := `INSERT INTO reports (id, title, project, company, author, version, status, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, report.ID, report.Title, report.Project, report.Company,
		report.Author, report.Version, report.Status, report.CreatedAt, report.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert report: %w", err)
	}
	return nil
}

func (r *ReportRepository) FindByID(id string) (*models.Report, error) {
	query := `SELECT id, title, project, company, author, version, status, COALESCE(pdf_path,''), created_at, updated_at
	          FROM reports WHERE id = ?`
	var m models.Report
	err := r.db.QueryRow(query, id).Scan(&m.ID, &m.Title, &m.Project, &m.Company,
		&m.Author, &m.Version, &m.Status, &m.PDFPath, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find report by id: %w", err)
	}
	return &m, nil
}

func (r *ReportRepository) UpdatePDFPath(id, pdfPath string) error {
	query := `UPDATE reports SET pdf_path = ?, status = 'completed', updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := r.db.Exec(query, pdfPath, id)
	if err != nil {
		return fmt.Errorf("update report pdf path: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("report not found")
	}
	return nil
}
