package response

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

type meta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type listResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    meta        `json:"meta"`
}

func JSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, APIResponse{Success: true, Data: data})
}

func List(w http.ResponseWriter, data interface{}, page, perPage, total int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listResponse{
		Success: true,
		Data:    data,
		Meta:    meta{Page: page, PerPage: perPage, Total: total},
	})
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, APIResponse{Success: true, Data: data})
}

func BadRequest(w http.ResponseWriter, err string) {
	JSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err})
}

func NotFound(w http.ResponseWriter, err string) {
	JSON(w, http.StatusNotFound, APIResponse{Success: false, Error: err})
}

func InternalError(w http.ResponseWriter, err string) {
	JSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err})
}
