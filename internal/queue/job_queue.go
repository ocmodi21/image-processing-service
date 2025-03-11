package queue

import (
	"sync"
)

// JobQueue represents a queue for processing jobs
type JobQueue struct {
	queue     []string
	mu        sync.Mutex
	cond      *sync.Cond
	processor func(string)
	workers   int
	stop      chan struct{}
	wg        sync.WaitGroup
}

// NewJobQueue creates a new job queue with the specified number of workers
func NewJobQueue(workers int, processor func(string)) *JobQueue {
	q := &JobQueue{
		queue:     make([]string, 0),
		workers:   workers,
		processor: processor,
		stop:      make(chan struct{}),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Start starts the job queue workers
func (q *JobQueue) Start() {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker()
	}
}

// Stop stops the job queue
func (q *JobQueue) Stop() {
	close(q.stop)
	q.cond.Broadcast()
	q.wg.Wait()
}

// Enqueue adds a job to the queue
func (q *JobQueue) Enqueue(jobID string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue = append(q.queue, jobID)
	q.cond.Signal()
}

// worker processes jobs from the queue
func (q *JobQueue) worker() {
	defer q.wg.Done()

	for {
		q.mu.Lock()
		for len(q.queue) == 0 {
			// Check if we should stop
			select {
			case <-q.stop:
				q.mu.Unlock()
				return
			default:
			}

			// Wait for a job
			q.cond.Wait()

			// Check again if we should stop
			select {
			case <-q.stop:
				q.mu.Unlock()
				return
			default:
			}
		}

		// Get the next job
		jobID := q.queue[0]
		q.queue = q.queue[1:]
		q.mu.Unlock()

		// Process the job
		q.processor(jobID)
	}
}
