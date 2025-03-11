package storage

import (
	"errors"
	"sync"

	"github.com/ocmodi21/image-processing-service/internal/models"
)

var (
	ErrJobNotFound = errors.New("job not found")
)

// JobStorage represents an in-memory storage for jobs
// In a production environment, this would be replaced with a database
type JobStorage struct {
	mu     sync.RWMutex
	jobs   map[string]*models.Job
	nextID int
}

func NewJobStorage() *JobStorage {
	return &JobStorage{
		jobs:   make(map[string]*models.Job),
		nextID: 1,
	}
}

func (s *JobStorage) CreateJob(job *models.Job) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate a job ID
	jobID := string(rune(s.nextID))
	s.nextID++

	job.ID = jobID
	s.jobs[jobID] = job

	return jobID
}

func (s *JobStorage) GetJob(jobID string) (*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return nil, ErrJobNotFound
	}

	return job, nil
}

func (s *JobStorage) UpdateJob(job *models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID]; !exists {
		return ErrJobNotFound
	}

	s.jobs[job.ID] = job
	return nil
}
