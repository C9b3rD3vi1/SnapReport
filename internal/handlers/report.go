package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/C9b3rD3vi1/SnapReport/internal/response"
	"github.com/C9b3rD3vi1/SnapReport/internal/service"
	"github.com/C9b3rD3vi1/SnapReport/internal/services/report"
)

type ReportHandler struct {
	service *report.Service
}

func NewReportHandler(svc *report.Service) *ReportHandler {
	return &ReportHandler{service: svc}
}

func (h *ReportHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req report.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Title == "" {
		response.BadRequest(w, "Report title is required")
		return
	}
	if len(req.Uploads) == 0 {
		response.BadRequest(w, "At least one upload is required")
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		if errors.Is(err, service.ErrTitleMissing) || errors.Is(err, service.ErrNoUploads) {
			response.BadRequest(w, err.Error())
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, result)
}

func (h *ReportHandler) List(w http.ResponseWriter, r *http.Request) {
	reports, err := h.service.List()
	if err != nil {
		response.InternalError(w, "Failed to list reports")
		return
	}
	response.OK(w, reports)
}

func (h *ReportHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Report ID is required")
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w, "Report not found")
			return
		}
		response.InternalError(w, "Failed to delete report")
		return
	}

	response.OK(w, nil)
}

func (h *ReportHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Report ID is required")
		return
	}

	reportData, uploads, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w, "Report not found")
			return
		}
		response.InternalError(w, "Failed to retrieve report")
		return
	}

	response.OK(w, map[string]interface{}{
		"report":  reportData,
		"uploads": uploads,
	})
}

func (h *ReportHandler) Preview(w http.ResponseWriter, r *http.Request) {
	var req report.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	if req.Title == "" {
		response.BadRequest(w, "Report title is required")
		return
	}

	html, err := h.service.PreviewHTML(req)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (h *ReportHandler) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Report ID is required")
		return
	}

	reportData, _, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w, "Report not found")
			return
		}
		response.InternalError(w, "Failed to retrieve report")
		return
	}

	if reportData.PDFPath == "" {
		response.NotFound(w, "PDF not yet generated")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, reportData.Title))
	http.ServeFile(w, r, reportData.PDFPath)
}
