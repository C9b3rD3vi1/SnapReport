package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	OK(w, map[string]string{"key": "value"}, "success")

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body APIResponse
	json.NewDecoder(resp.Body).Decode(&body)

	if !body.Success {
		t.Error("expected success true")
	}
	if body.Message != "success" {
		t.Errorf("expected message 'success', got %q", body.Message)
	}
	if body.Error != "" {
		t.Errorf("expected no error, got %q", body.Error)
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	Created(w, "data", "created")

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var body APIResponse
	json.NewDecoder(resp.Body).Decode(&body)

	if !body.Success {
		t.Error("expected success true")
	}
	if body.Data != "data" {
		t.Errorf("expected data 'data', got %v", body.Data)
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequest(w, "invalid input")

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body APIResponse
	json.NewDecoder(resp.Body).Decode(&body)

	if body.Success {
		t.Error("expected success false")
	}
	if body.Error != "invalid input" {
		t.Errorf("expected error 'invalid input', got %q", body.Error)
	}
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	NotFound(w, "not found")

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	InternalError(w, "server error")

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusTeapot, APIResponse{Success: true, Message: "teapot"})

	resp := w.Result()
	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("expected status 418, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}
