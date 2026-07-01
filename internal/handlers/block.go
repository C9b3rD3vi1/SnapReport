package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/C9b3rD3vi1/SnapReport/internal/fileutil"
	"github.com/C9b3rD3vi1/SnapReport/internal/models"
	"github.com/C9b3rD3vi1/SnapReport/internal/repository"
	"github.com/C9b3rD3vi1/SnapReport/internal/response"
)

type BlockHandler struct {
	repo *repository.BlockRepository
}

func NewBlockHandler(repo *repository.BlockRepository) *BlockHandler {
	return &BlockHandler{repo: repo}
}

type createBlockRequest struct {
	Type     models.BlockType `json:"type"`
	Position int              `json:"position"`
	Content  string           `json:"content"`
}

func (h *BlockHandler) Create(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "reportID")
	if reportID == "" {
		response.BadRequest(w, "Report ID required")
		return
	}

	var req createBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid body")
		return
	}

	block := &models.Block{
		ID:        fileutil.GenerateID(),
		ReportID:  reportID,
		Type:      req.Type,
		Position:  req.Position,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.Insert(block); err != nil {
		response.InternalError(w, "Failed to create block")
		return
	}
	response.Created(w, block)
}

func (h *BlockHandler) List(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "reportID")
	if reportID == "" {
		response.BadRequest(w, "Report ID required")
		return
	}

	blocks, err := h.repo.FindByReportID(reportID)
	if err != nil {
		response.InternalError(w, "Failed to list blocks")
		return
	}
	if blocks == nil {
		blocks = []models.Block{}
	}
	response.OK(w, blocks)
}

func (h *BlockHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Block ID required")
		return
	}

	existing, err := h.repo.FindByID(id)
	if err != nil {
		response.InternalError(w, "Failed to find block")
		return
	}
	if existing == nil {
		response.NotFound(w, "Block not found")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		response.BadRequest(w, "Invalid body")
		return
	}

	if v, ok := updates["content"].(string); ok {
		existing.Content = v
	}
	if v, ok := updates["position"].(float64); ok {
		existing.Position = int(v)
	}
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(existing); err != nil {
		response.InternalError(w, "Failed to update block")
		return
	}
	response.OK(w, existing)
}

func (h *BlockHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "Block ID required")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		response.InternalError(w, "Failed to delete block")
		return
	}
	response.OK(w, nil)
}

type reorderRequest struct {
	BlockIDs []string `json:"block_ids"`
}

func (h *BlockHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "reportID")
	if reportID == "" {
		response.BadRequest(w, "Report ID required")
		return
	}

	var req reorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid body")
		return
	}

	if err := h.repo.Reorder(reportID, req.BlockIDs); err != nil {
		response.InternalError(w, "Failed to reorder")
		return
	}
	response.OK(w, nil)
}
