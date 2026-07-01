package service

import "github.com/C9b3rD3vi1/SnapReport/internal/models"

type UploadRepo interface {
	Insert(*models.Upload) error
	FindAll() ([]models.Upload, error)
	FindByID(string) (*models.Upload, error)
	FindByReportID(string) ([]models.Upload, error)
	UpdateMetadata(id, title, description, notes string, orderIndex int, reportID string) error
	Delete(id string) error
}

type ReportRepo interface {
	Insert(*models.Report) error
	FindByID(string) (*models.Report, error)
	FindAll() ([]models.Report, error)
	Delete(id string) error
	UpdatePDFPath(id, path string) error
}
