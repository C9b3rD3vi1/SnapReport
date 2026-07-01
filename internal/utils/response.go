package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func JSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func OK(w http.ResponseWriter, data interface{}, message string) {
	JSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Created(w http.ResponseWriter, data interface{}, message string) {
	JSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func BadRequest(w http.ResponseWriter, err string) {
	JSON(w, http.StatusBadRequest, APIResponse{
		Success: false,
		Error:   err,
	})
}

func NotFound(w http.ResponseWriter, err string) {
	JSON(w, http.StatusNotFound, APIResponse{
		Success: false,
		Error:   err,
	})
}

func InternalError(w http.ResponseWriter, err string) {
	JSON(w, http.StatusInternalServerError, APIResponse{
		Success: false,
		Error:   err,
	})
}
