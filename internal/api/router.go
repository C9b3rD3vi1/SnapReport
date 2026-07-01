package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/anomalyco/SnapReport/internal/handlers"
	"github.com/anomalyco/SnapReport/internal/middleware"
)

type Handlers struct {
	Upload *handlers.UploadHandler
	Report *handlers.ReportHandler
}

func NewRouter(h *Handlers, uploadDir string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logging)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Get("/health", handlers.HealthCheck)

	fileServer := http.FileServer(http.Dir(uploadDir))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/uploads", func(r chi.Router) {
			r.Post("/", h.Upload.Create)
			r.Get("/", h.Upload.List)
			r.Delete("/{id}", h.Upload.Delete)
		})

		r.Route("/reports", func(r chi.Router) {
			r.Post("/", h.Report.Create)
			r.Get("/{id}", h.Report.Get)
			r.Get("/{id}/download", h.Report.Download)
		})
	})

	return r
}
