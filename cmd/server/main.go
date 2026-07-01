package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/C9b3rD3vi1/SnapReport/internal/api"
	"github.com/C9b3rD3vi1/SnapReport/internal/config"
	"github.com/C9b3rD3vi1/SnapReport/internal/handlers"
	"github.com/C9b3rD3vi1/SnapReport/internal/middleware"
	"github.com/C9b3rD3vi1/SnapReport/internal/pdf"
	"github.com/C9b3rD3vi1/SnapReport/internal/repository"
	reportSvc "github.com/C9b3rD3vi1/SnapReport/internal/services/report"
	"github.com/C9b3rD3vi1/SnapReport/internal/services"
	"github.com/C9b3rD3vi1/SnapReport/internal/service"
)

func main() {
	cfg := config.Load()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	db, err := repository.NewSQLite(cfg.DatabasePath)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	dirs := []string{cfg.UploadDir, cfg.PDFDir}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			slog.Error("failed to create directory", "dir", d, "error", err)
			os.Exit(1)
		}
	}

	var uploadRepo service.UploadRepo = repository.NewUploadRepository(db.DB())
	var reportRepo service.ReportRepo = repository.NewReportRepository(db.DB())
	blockRepo := repository.NewBlockRepository(db.DB())

	uploadSvc := services.NewUploadService(uploadRepo, cfg.UploadDir, cfg.MaxUploadMB)

	pdfGen, err := pdf.NewGenerator()
	if err != nil {
		slog.Warn("pdf generator initialization failed, pdf generation disabled", "error", err)
		pdfGen = nil
	}

	reportSvc := reportSvc.NewService(reportRepo, uploadRepo, blockRepo, pdfGen, cfg.PDFDir)

	uploadHandler := handlers.NewUploadHandler(uploadSvc)
	reportHandler := handlers.NewReportHandler(reportSvc)
	templateHandler := handlers.NewTemplateHandler()
	blockHandler := handlers.NewBlockHandler(blockRepo)

	h := &api.Handlers{
		Upload:   uploadHandler,
		Report:   reportHandler,
		Template: templateHandler,
		Block:    blockHandler,
	}

	rateLimiter := middleware.NewRateLimiter(100, 1*time.Minute)
	router := api.NewRouter(h, cfg.UploadDir, cfg.AllowedOrigins, rateLimiter)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		slog.Info("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	slog.Info("server starting",
		"port", cfg.Port,
		"environment", cfg.Environment,
	)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
