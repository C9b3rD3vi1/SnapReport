package repository

import (
	"database/sql"
	"fmt"
	"github.com/C9b3rD3vi1/SnapReport/internal/models"
)

type BlockRepository struct {
	db *sql.DB
}

func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{db: db}
}

func (r *BlockRepository) Insert(b *models.Block) error {
	_, err := r.db.Exec(
		`INSERT INTO blocks (id, report_id, type, position, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		b.ID, b.ReportID, b.Type, b.Position, b.Content, b.CreatedAt, b.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert block: %w", err)
	}
	return nil
}

func (r *BlockRepository) FindByReportID(reportID string) ([]models.Block, error) {
	rows, err := r.db.Query(
		`SELECT id, report_id, type, position, content, created_at, updated_at FROM blocks WHERE report_id = ? ORDER BY position ASC`,
		reportID,
	)
	if err != nil {
		return nil, fmt.Errorf("find blocks: %w", err)
	}
	defer rows.Close()

	var blocks []models.Block
	for rows.Next() {
		var b models.Block
		if err := rows.Scan(&b.ID, &b.ReportID, &b.Type, &b.Position, &b.Content, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan block: %w", err)
		}
		blocks = append(blocks, b)
	}
	return blocks, rows.Err()
}

func (r *BlockRepository) Update(b *models.Block) error {
	_, err := r.db.Exec(
		`UPDATE blocks SET type=?, position=?, content=?, updated_at=? WHERE id=?`,
		b.Type, b.Position, b.Content, b.UpdatedAt, b.ID,
	)
	if err != nil {
		return fmt.Errorf("update block: %w", err)
	}
	return nil
}

func (r *BlockRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM blocks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete block: %w", err)
	}
	return nil
}

func (r *BlockRepository) DeleteByReport(reportID string) error {
	_, err := r.db.Exec(`DELETE FROM blocks WHERE report_id = ?`, reportID)
	return err
}

func (r *BlockRepository) Reorder(reportID string, blockIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	for i, id := range blockIDs {
		if _, err := tx.Exec(`UPDATE blocks SET position=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND report_id=?`, i, id, reportID); err != nil {
			return fmt.Errorf("reorder block: %w", err)
		}
	}
	return tx.Commit()
}

func (r *BlockRepository) FindByID(id string) (*models.Block, error) {
	var b models.Block
	err := r.db.QueryRow(
		`SELECT id, report_id, type, position, content, created_at, updated_at FROM blocks WHERE id = ?`, id,
	).Scan(&b.ID, &b.ReportID, &b.Type, &b.Position, &b.Content, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find block: %w", err)
	}
	return &b, nil
}
