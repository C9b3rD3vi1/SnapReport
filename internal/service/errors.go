package service

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidFile  = errors.New("invalid file")
	ErrFileTooLarge = errors.New("file too large")
	ErrTitleMissing = errors.New("report title is required")
	ErrNoUploads    = errors.New("at least one upload is required")
)
