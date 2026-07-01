package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/anomalyco/SnapReport/internal/services/report"
	"github.com/anomalyco/SnapReport/internal/utils"
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
		utils.BadRequest(w, "Invalid request body")
		return
	}

	if req.Title == "" {
		utils.BadRequest(w, "Report title is required")
		return
	}

	if len(req.Uploads) == 0 {
		utils.BadRequest(w, "At least one upload is required")
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.Created(w, result, "Report created successfully")
}

func (h *ReportHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.BadRequest(w, "Report ID is required")
		return
	}

	reportData, uploads, err := h.service.Get(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFound(w, "Report not found")
			return
		}
		utils.InternalError(w, "Failed to retrieve report")
		return
	}

	utils.OK(w, map[string]interface{}{
		"report":  reportData,
		"uploads": uploads,
	}, "Report retrieved successfully")
}

func (h *ReportHandler) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.BadRequest(w, "Report ID is required")
		return
	}

	reportData, _, err := h.service.Get(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.NotFound(w, "Report not found")
			return
		}
		utils.InternalError(w, "Failed to retrieve report")
		return
	}

	if reportData.PDFPath == "" {
		utils.NotFound(w, "PDF not yet generated")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, reportData.Title))
	http.ServeFile(w, r, reportData.PDFPath)
}
