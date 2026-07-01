package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/C9b3rD3vi1/SnapReport/internal/response"
	"github.com/C9b3rD3vi1/SnapReport/internal/templates"
)

type TemplateHandler struct{}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{}
}

func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	response.OK(w, templates.All)
}

func (h *TemplateHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	t := templates.FindByID(id)
	if t == nil {
		response.NotFound(w, "Template not found")
		return
	}
	response.OK(w, t)
}
