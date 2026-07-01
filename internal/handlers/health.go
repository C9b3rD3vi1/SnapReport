package handlers

import (
	"net/http"
	"time"

	"github.com/C9b3rD3vi1/SnapReport/internal/response"
)

var startTime = time.Now()

type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Uptime    string `json:"uptime"`
	Timestamp string `json:"timestamp"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.OK(w, HealthResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).Round(time.Second).String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
