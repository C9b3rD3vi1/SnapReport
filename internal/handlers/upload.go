package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/anomalyco/SnapReport/internal/services"
	"github.com/anomalyco/SnapReport/internal/utils"
)

type UploadHandler struct {
	service *services.UploadService
}

func NewUploadHandler(svc *services.UploadService) *UploadHandler {
	return &UploadHandler{service: svc}
}

func (h *UploadHandler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		utils.BadRequest(w, "Request must be multipart/form-data")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		utils.BadRequest(w, "Failed to parse upload form")
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		utils.BadRequest(w, "No files provided in 'files' field")
		return
	}

	uploads, err := h.service.Create(files)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.Created(w, uploads, "Files uploaded successfully")
}

func (h *UploadHandler) List(w http.ResponseWriter, r *http.Request) {
	uploads, err := h.service.List()
	if err != nil {
		utils.InternalError(w, "Failed to list uploads")
		return
	}

	utils.OK(w, uploads, "Uploads retrieved successfully")
}

func (h *UploadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.BadRequest(w, "Upload ID is required")
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, errors.New("upload not found")) {
			utils.NotFound(w, "Upload not found")
			return
		}
		utils.InternalError(w, "Failed to delete upload")
		return
	}

	utils.OK(w, nil, "Upload deleted successfully")
}
