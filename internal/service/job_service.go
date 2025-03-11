package service

import (
	"errors"
	"time"

	"github.com/ocmodi21/image-processing-service/internal/models"
	"github.com/ocmodi21/image-processing-service/internal/queue"
	"github.com/ocmodi21/image-processing-service/internal/storage"
)

var (
	ErrInvalidRequest = errors.New("invalid request: count does not match number of visits")
)

// JobService handles the business logic for job management
type JobService struct {
	jobStorage     *storage.JobStorage
	storeStorage   *storage.StoreStorage
	jobQueue       *queue.JobQueue
	imageProcessor *ImageProcessor
}

func NewJobService(
	jobStorage *storage.JobStorage,
	storeStorage *storage.StoreStorage,
	jobQueue *queue.JobQueue,
	imageProcessor *ImageProcessor,
) *JobService {
	return &JobService{
		jobStorage:     jobStorage,
		storeStorage:   storeStorage,
		jobQueue:       jobQueue,
		imageProcessor: imageProcessor,
	}
}

// CreateJob creates a new job and enqueues it for processing
func (s *JobService) CreateJob(req *models.JobSubmissionRequest) (string, error) {
	// Validate request
	if req.Count != len(req.Visits) {
		return "", ErrInvalidRequest
	}

	// Create a new job
	job := &models.Job{
		Status:    models.JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Visits:    req.Visits,
		Count:     req.Count,
	}

	// Store the job
	jobID := s.jobStorage.CreateJob(job)

	// Enqueue the job for processing
	s.jobQueue.Enqueue(jobID)

	return jobID, nil
}

// GetJob retrieves a job by ID
func (s *JobService) GetJob(jobID string) (*models.Job, error) {
	return s.jobStorage.GetJob(jobID)
}

// ProcessJob processes a job
func (s *JobService) ProcessJob(jobID string) {
	job, err := s.jobStorage.GetJob(jobID)
	if err != nil {
		// Log error
		return
	}

	// Update job status to ongoing
	job.Status = models.JobStatusOngoing
	job.UpdatedAt = time.Now()
	s.jobStorage.UpdateJob(job)

	// Process each visit
	hasErrors := false
	for i, visit := range job.Visits {
		// Check if store exists
		_, err := s.storeStorage.GetStore(visit.StoreID)
		if err != nil {
			hasErrors = true
			job.Errors = append(job.Errors, models.JobError{
				StoreID: visit.StoreID,
				Error:   "store not found",
			})
			continue
		}

		// Process images for this visit
		results, err := s.imageProcessor.ProcessImages(visit.ImageURLs)
		if err != nil {
			hasErrors = true
			job.Errors = append(job.Errors, models.JobError{
				StoreID: visit.StoreID,
				Error:   err.Error(),
			})
			continue
		}

		// Update visit with results
		job.Visits[i].Results = results
	}

	// Update job status based on processing results
	if hasErrors {
		job.Status = models.JobStatusFailed
	} else {
		job.Status = models.JobStatusCompleted
	}

	job.UpdatedAt = time.Now()
	s.jobStorage.UpdateJob(job)
}
