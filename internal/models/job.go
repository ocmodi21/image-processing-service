package models

import (
	"time"
)

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusOngoing   JobStatus = "ongoing"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type Job struct {
	ID        string     `json:"job_id"`
	Status    JobStatus  `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Visits    []Visit    `json:"visits"`
	Count     int        `json:"count"`
	Errors    []JobError `json:"errors,omitempty"`
}

type Visit struct {
	StoreID   string        `json:"store_id"`
	ImageURLs []string      `json:"image_url"`
	VisitTime string        `json:"visit_time"`
	Results   []ImageResult `json:"results,omitempty"`
}

type ImageResult struct {
	URL         string    `json:"url"`
	Perimeter   float64   `json:"perimeter"`
	ProcessedAt time.Time `json:"processed_at"`
}

type JobError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

type JobSubmissionRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

type JobSubmissionResponse struct {
	JobID string `json:"job_id"`
}

type JobStatusResponse struct {
	Status JobStatus  `json:"status"`
	JobID  string     `json:"job_id"`
	Errors []JobError `json:"error,omitempty"`
}
