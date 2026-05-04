package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ashczar77/mockingjay/backend/models"
	"github.com/ashczar77/mockingjay/backend/repository"
)

// Handler holds dependencies for all HTTP handlers
type Handler struct {
	repo repository.Repository
}

// New creates a new Handler
func New(repo repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateResult(w http.ResponseWriter, r *http.Request) {
	var res models.TestResult
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.repo.CreateResult(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, created)
}

func (h *Handler) GetResults(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}
	results, err := h.repo.GetResults(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, results)
}

func (h *Handler) CreateABTest(w http.ResponseWriter, r *http.Request) {
	var t models.ABTestResult
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.repo.CreateABTest(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, created)
}

func (h *Handler) GetABTests(w http.ResponseWriter, r *http.Request) {
	tests, err := h.repo.GetABTests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, tests)
}

func (h *Handler) CreateTranscription(w http.ResponseWriter, r *http.Request) {
	var t models.Transcription
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.repo.CreateTranscription(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, created)
}

func (h *Handler) GetTranscriptions(w http.ResponseWriter, r *http.Request) {
	transcriptions, err := h.repo.GetTranscriptions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, transcriptions)
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.repo.GetMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, metrics)
}

func (h *Handler) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.repo.GetHealthStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, status)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.Ping(); err != nil {
		http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(w, "OK")
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
