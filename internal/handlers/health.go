package handlers

import (
	"net/http"
	"github.com/anomalyco/SnapReport/internal/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.OK(w, map[string]string{
		"status": "healthy",
	}, "Server is running")
}
