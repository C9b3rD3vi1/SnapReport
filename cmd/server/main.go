package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anomalyco/SnapReport/internal/api"
	"github.com/anomalyco/SnapReport/internal/config"
	"github.com/anomalyco/SnapReport/internal/handlers"
	"github.com/anomalyco/SnapReport/internal/pdf"
	"github.com/anomalyco/SnapReport/internal/repository"
	reportSvc "github.com/anomalyco/SnapReport/internal/services/report"
	"github.com/anomalyco/SnapReport/internal/services"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

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

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		slog.Error("failed to create upload directory", "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.PDFDir, 0755); err != nil {
		slog.Error("failed to create PDF directory", "error", err)
		os.Exit(1)
	}

	uploadRepo := repository.NewUploadRepository(db.DB())
	reportRepo := repository.NewReportRepository(db.DB())

	uploadSvc := services.NewUploadService(uploadRepo, cfg.UploadDir, cfg.MaxUploadMB)

	pdfGen, err := pdf.NewGenerator()
	if err != nil {
		slog.Warn("pdf generator initialization failed, pdf generation disabled", "error", err)
		pdfGen = nil
	}

	reportSvc := reportSvc.NewService(reportRepo, uploadRepo, pdfGen, cfg.PDFDir)

	uploadHandler := handlers.NewUploadHandler(uploadSvc)
	reportHandler := handlers.NewReportHandler(reportSvc)

	h := &api.Handlers{
		Upload: uploadHandler,
		Report: reportHandler,
	}

	router := api.NewRouter(h, cfg.UploadDir)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		slog.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("server shutdown failed", "error", err)
		}
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
