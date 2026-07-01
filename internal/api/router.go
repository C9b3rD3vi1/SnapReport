package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/C9b3rD3vi1/SnapReport/internal/handlers"
	"github.com/C9b3rD3vi1/SnapReport/internal/middleware"
)

type Handlers struct {
	Upload   *handlers.UploadHandler
	Report   *handlers.ReportHandler
	Template *handlers.TemplateHandler
	Block    *handlers.BlockHandler
}

func NewRouter(h *Handlers, uploadDir string, allowedOrigins []string, rl *middleware.RateLimiter) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logging)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS(allowedOrigins))
	r.Use(middleware.SecurityHeaders)
	r.Use(rl.Middleware)

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
			r.Get("/", h.Report.List)
			r.Post("/", h.Report.Create)
			r.Post("/preview", h.Report.Preview)
			r.Get("/{id}", h.Report.Get)
			r.Get("/{id}/download", h.Report.Download)
			r.Delete("/{id}", h.Report.Delete)

			r.Route("/{reportID}/blocks", func(r chi.Router) {
				r.Get("/", h.Block.List)
				r.Post("/", h.Block.Create)
				r.Post("/reorder", h.Block.Reorder)
				r.Patch("/{id}", h.Block.Update)
				r.Delete("/{id}", h.Block.Delete)
			})
		})

		r.Route("/templates", func(r chi.Router) {
			r.Get("/", h.Template.List)
			r.Get("/{id}", h.Template.Get)
		})
	})

	return r
}
