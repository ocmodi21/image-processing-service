package api

import (
	"encoding/json"
	"net/http"

	"github.com/ocmodi21/image-processing-service/internal/models"
	"github.com/ocmodi21/image-processing-service/internal/service"
	"github.com/ocmodi21/image-processing-service/internal/storage"
)

// Handler handles HTTP requests
type Handler struct {
	jobService *service.JobService
}

func NewHandler(jobService *service.JobService) *Handler {
	return &Handler{
		jobService: jobService,
	}
}

// SubmitJob handles job submission requests
func (h *Handler) SubmitJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.JobSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	jobID, err := h.jobService.CreateJob(&req)
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			respondWithError(w, http.StatusBadRequest, "Invalid request: count does not match number of visits")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to create job")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, models.JobSubmissionResponse{JobID: jobID})
}

// GetJobStatus handles job status requests
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobID := r.URL.Query().Get("jobid")
	if jobID == "" {
		respondWithError(w, http.StatusBadRequest, "Job ID is required")
		return
	}

	job, err := h.jobService.GetJob(jobID)
	if err != nil {
		switch err {
		case storage.ErrJobNotFound:
			respondWithJSON(w, http.StatusBadRequest, map[string]string{})
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to get job status")
		}
		return
	}

	response := models.JobStatusResponse{
		Status: job.Status,
		JobID:  job.ID,
	}

	if job.Status == models.JobStatusFailed {
		response.Errors = job.Errors
	}

	respondWithJSON(w, http.StatusOK, response)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
