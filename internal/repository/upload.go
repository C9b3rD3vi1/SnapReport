package repository

import (
	"database/sql"
	"fmt"
	"github.com/C9b3rD3vi1/SnapReport/internal/models"
)

type UploadRepository struct {
	db *sql.DB
}

func NewUploadRepository(db *sql.DB) *UploadRepository {
	return &UploadRepository{db: db}
}

var uploadCols = "id, filename, original_name, mime_type, size, path, " +
	"COALESCE(title,''), COALESCE(description,''), COALESCE(notes,''), " +
	"COALESCE(category,''), COALESCE(priority,''), COALESCE(severity,''), COALESCE(status,''), COALESCE(recommendation,''), " +
	"order_index, created_at"

func scanUpload(scanner interface {
	Scan(dest ...interface{}) error
}) (models.Upload, error) {
	var u models.Upload
	err := scanner.Scan(&u.ID, &u.Filename, &u.OriginalName, &u.MimeType, &u.Size, &u.Path,
		&u.Title, &u.Description, &u.Notes,
		&u.Category, &u.Priority, &u.Severity, &u.Status, &u.Recommendation,
		&u.OrderIndex, &u.CreatedAt)
	return u, err
}

func (r *UploadRepository) Insert(u *models.Upload) error {
	query := `INSERT INTO uploads (id, filename, original_name, mime_type, size, path, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.ID, u.Filename, u.OriginalName, u.MimeType, u.Size, u.Path, u.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert upload: %w", err)
	}
	return nil
}

func (r *UploadRepository) FindAll() ([]models.Upload, error) {
	query := `SELECT ` + uploadCols + ` FROM uploads ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("find all uploads: %w", err)
	}
	defer rows.Close()

	var uploads []models.Upload
	for rows.Next() {
		u, err := scanUpload(rows)
		if err != nil {
			return nil, fmt.Errorf("scan upload: %w", err)
		}
		uploads = append(uploads, u)
	}
	return uploads, rows.Err()
}

func (r *UploadRepository) FindByID(id string) (*models.Upload, error) {
	query := `SELECT ` + uploadCols + ` FROM uploads WHERE id = ?`
	row := r.db.QueryRow(query, id)
	u, err := scanUpload(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find upload by id: %w", err)
	}
	return &u, nil
}

func (r *UploadRepository) UpdateMetadata(id, title, description, notes, category, priority, severity, status, recommendation string, orderIndex int, reportID string) error {
	query := `UPDATE uploads SET title=?, description=?, notes=?, category=?, priority=?, severity=?, status=?, recommendation=?, order_index=?, report_id=? WHERE id=?`
	result, err := r.db.Exec(query, title, description, notes, category, priority, severity, status, recommendation, orderIndex, reportID, id)
	if err != nil {
		return fmt.Errorf("update upload metadata: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("upload not found")
	}
	return nil
}

func (r *UploadRepository) FindByReportID(reportID string) ([]models.Upload, error) {
	query := `SELECT ` + uploadCols + ` FROM uploads WHERE report_id = ? ORDER BY order_index ASC`
	rows, err := r.db.Query(query, reportID)
	if err != nil {
		return nil, fmt.Errorf("find uploads by report id: %w", err)
	}
	defer rows.Close()

	var uploads []models.Upload
	for rows.Next() {
		u, err := scanUpload(rows)
		if err != nil {
			return nil, fmt.Errorf("scan upload: %w", err)
		}
		uploads = append(uploads, u)
	}
	return uploads, rows.Err()
}

func (r *UploadRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM uploads WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete upload: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("upload not found")
	}
	return nil
}
