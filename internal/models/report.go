package models

import "time"

type Report struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Project     string    `json:"project"`
	Company     string    `json:"company"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	Status      string    `json:"status"`
	PDFPath     string    `json:"pdf_path,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ReportInfo struct {
	Title   string `json:"title"`
	Project string `json:"project"`
	Company string `json:"company"`
	Author  string `json:"author"`
	Version string `json:"version"`
}
