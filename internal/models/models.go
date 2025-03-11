// internal/models/models.go
package models

import (
	"time"
)

// JobStatus represents the current state of a job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

// Visit represents a store visit with images
type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

// SubmitJobRequest is the payload for job submission
type SubmitJobRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

// SubmitJobResponse is the response after job submission
type SubmitJobResponse struct {
	JobID string `json:"job_id"`
}

// JobStatusResponse is the response for job status queries
type JobStatusResponse struct {
	Status JobStatus  `json:"status"`
	JobID  string     `json:"job_id"`
	Errors []JobError `json:"error,omitempty"`
}

// JobError represents an error that occurred during job processing
type JobError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

// ImageResult stores the processing result for a single image
type ImageResult struct {
	StoreID   string
	ImageURL  string
	Perimeter float64
	Error     string
}

// Job represents a complete processing job
type Job struct {
	ID        string
	Status    JobStatus
	Visits    []Visit
	Results   []ImageResult
	Errors    []JobError
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Store represents store information from the master data
type Store struct {
	ID       string
	Name     string
	AreaCode string
}
