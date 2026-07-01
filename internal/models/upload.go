package models

import (
	"path/filepath"
	"strings"
	"time"
)

type Upload struct {
	ID           string    `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	Path         string    `json:"-"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Notes        string    `json:"notes,omitempty"`
	OrderIndex   int       `json:"order_index"`
	ReportID     string    `json:"report_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
}

func (u *Upload) SetThumbnailURL() {
	ext := filepath.Ext(u.Filename)
	base := strings.TrimSuffix(u.Filename, ext)
	u.ThumbnailURL = "/uploads/thumbnails/" + base + "_thumb.jpg"
}

type UploadMetadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Notes       string `json:"notes,omitempty"`
	OrderIndex  int    `json:"order_index"`
}
