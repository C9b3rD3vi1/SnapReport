package repository

import (
	"database/sql"
	"fmt"
	"github.com/anomalyco/SnapReport/internal/models"
)

type UploadRepository struct {
	db *sql.DB
}

func NewUploadRepository(db *sql.DB) *UploadRepository {
	return &UploadRepository{db: db}
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
	query := `SELECT id, filename, original_name, mime_type, size, path, COALESCE(title,''), COALESCE(description,''), COALESCE(notes,''), order_index, created_at
	          FROM uploads ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("find all uploads: %w", err)
	}
	defer rows.Close()

	var uploads []models.Upload
	for rows.Next() {
		var u models.Upload
		if err := rows.Scan(&u.ID, &u.Filename, &u.OriginalName, &u.MimeType, &u.Size, &u.Path,
			&u.Title, &u.Description, &u.Notes, &u.OrderIndex, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan upload: %w", err)
		}
		uploads = append(uploads, u)
	}
	return uploads, rows.Err()
}

func (r *UploadRepository) FindByID(id string) (*models.Upload, error) {
	query := `SELECT id, filename, original_name, mime_type, size, path, COALESCE(title,''), COALESCE(description,''), COALESCE(notes,''), order_index, created_at
	          FROM uploads WHERE id = ?`
	var u models.Upload
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Filename, &u.OriginalName, &u.MimeType, &u.Size, &u.Path,
		&u.Title, &u.Description, &u.Notes, &u.OrderIndex, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find upload by id: %w", err)
	}
	return &u, nil
}

func (r *UploadRepository) UpdateMetadata(id, title, description, notes string, orderIndex int, reportID string) error {
	query := `UPDATE uploads SET title = ?, description = ?, notes = ?, order_index = ?, report_id = ? WHERE id = ?`
	result, err := r.db.Exec(query, title, description, notes, orderIndex, reportID, id)
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
	query := `SELECT id, filename, original_name, mime_type, size, path, COALESCE(title,''), COALESCE(description,''), COALESCE(notes,''), order_index, created_at
	          FROM uploads WHERE report_id = ? ORDER BY order_index ASC`
	rows, err := r.db.Query(query, reportID)
	if err != nil {
		return nil, fmt.Errorf("find uploads by report id: %w", err)
	}
	defer rows.Close()

	var uploads []models.Upload
	for rows.Next() {
		var u models.Upload
		if err := rows.Scan(&u.ID, &u.Filename, &u.OriginalName, &u.MimeType, &u.Size, &u.Path,
			&u.Title, &u.Description, &u.Notes, &u.OrderIndex, &u.CreatedAt); err != nil {
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
