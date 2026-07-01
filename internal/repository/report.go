package repository

import (
	"database/sql"
	"fmt"
	"github.com/C9b3rD3vi1/SnapReport/internal/models"
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

func (r *ReportRepository) FindAll() ([]models.Report, error) {
	query := `SELECT id, title, project, company, author, version, status, COALESCE(pdf_path,''), created_at, updated_at
	          FROM reports ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("find all reports: %w", err)
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var m models.Report
		if err := rows.Scan(&m.ID, &m.Title, &m.Project, &m.Company,
			&m.Author, &m.Version, &m.Status, &m.PDFPath, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		reports = append(reports, m)
	}
	return reports, rows.Err()
}

func (r *ReportRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM reports WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete report: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("report not found")
	}
	return nil
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
