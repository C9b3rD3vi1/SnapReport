package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/C9b3rD3vi1/SnapReport/internal/response"
	"github.com/C9b3rD3vi1/SnapReport/internal/service"
	"github.com/C9b3rD3vi1/SnapReport/internal/services"
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
		response.BadRequest(w, "Request must be multipart/form-data")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.BadRequest(w, "Failed to parse upload form")
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		response.BadRequest(w, "No files provided in 'files' field")
		return
	}

	uploads, errs := h.service.Create(files)
	if len(errs) > 0 && len(uploads) == 0 {
		response.BadRequest(w, errs[0].Err.Error())
		return
	}
	if len(errs) > 0 {
		response.Created(w, map[string]interface{}{
			"uploads": uploads,
			"errors":  errs,
		})
		return
	}

	response.Created(w, uploads)
}

func (h *UploadHandler) List(w http.ResponseWriter, r *http.Request) {
	uploads, err := h.service.List()
	if err != nil {
		response.InternalError(w, "Failed to list uploads")
		return
	}
	response.OK(w, uploads)
}

func (h *UploadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Upload ID is required")
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w, "Upload not found")
			return
		}
		response.InternalError(w, "Failed to delete upload")
		return
	}

	response.OK(w, nil)
}
