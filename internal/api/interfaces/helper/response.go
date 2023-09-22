package helper

import (
	"encoding/json"
	"net/http"
)

// response implements JSON response payload structure
type response struct {
	Status string          `json:"status"`
	Result json.RawMessage `json:"result,omitempty"`
}

// responseError implements Error structure to return in response payloads.
type responseError struct {
	ID      string   `json:"id"`
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

// ErrorRes implements JSON error response payload structure
type errorRes struct {
	Status string        `json:"status"`
	Error  responseError `json:"error"`
}

// Fail ends an unsuccessful JSON response
func Fail(w http.ResponseWriter, errorCode int) {

	// no cache
	w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")
	w.Header().Set("Expires", "-1")
	w.Header().Set("Pragma", "no-cache")

	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	w.WriteHeader(errorCode)
}

// Succeed sends a successful JSON response (200)
func Succeed(w http.ResponseWriter, result interface{}) {
	rj, err := json.Marshal(result)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	s := success(rj)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(s)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// Created sends a successful JSON with created new record response (201)
func Created(w http.ResponseWriter, result interface{}) {
	rj, err := json.Marshal(result)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	s := success(rj)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(s)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// NoContent sends a successful response (204)
func NoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	w.WriteHeader(http.StatusNoContent)
}

func success(res json.RawMessage) *response {
	return &response{
		Status: "ok",
		Result: res,
	}
}
