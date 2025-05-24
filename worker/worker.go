// main/worker/worker.go
package worker

import (
	"context"
	"fmt"
	"github/Eranmonnie/jqueue/jobs"
	"github/Eranmonnie/jqueue/queue"
	"log"
	"math"
	"sync"
	"time"
)

type JobHandler func(payload interface{}) error

// WorkerRegistry stores job handlers with thread-safe access
type WorkerRegistry struct {
	handlers map[string]JobHandler
	mutex    sync.RWMutex
}

// Global registry instance
var Registry = &WorkerRegistry{
	handlers: make(map[string]JobHandler),
}

func (wr *WorkerRegistry) RegisterHandler(jobName string, handler JobHandler) error {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	wr.handlers[jobName] = handler
	log.Printf("Registered handler for job type: %s", jobName)
	return nil

	//put the worker in the db just for the sake of it
}

func (wr *WorkerRegistry) GetHandler(jobType string) (JobHandler, error) {
	wr.mutex.RLock()
	defer wr.mutex.RUnlock()

	handler, exists := wr.handlers[jobType]
	if !exists {
		return nil, fmt.Errorf("no handler registered for job type: %s", jobType)
	}

	return handler, nil
}

// Worker function
func StartWorkerPool(workerCount int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := range workerCount {
		go runWorker(ctx, i)
	}
}

func runWorker(ctx context.Context, id int) {
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		default:
			// Process a job
			job, err := queue.DequeueJob()
			if err != nil {
				log.Printf("Worker %d error: %v", id, err)
				continue
			}

			processJobWithRetry(job)
		}
	}
}

func processJobWithRetry(job jobs.Job) {
	retries := 0
	maxRetries := job.MaxRetries

	for retries <= maxRetries {
		err := processJob(job)
		if err == nil {
			// Job succeeded
			log.Printf("Job %d completed successfully", job.ID)
			// database.MarkJobCompleted(job.ID, nil)
			return
		}

		retries++
		if retries > maxRetries {
			// Give up after max retries
			log.Printf("Job %d failed after %d retries: %v", job.ID, retries, err)
			// database.MarkJobFailed(job.ID, []byte(fmt.Sprintf(`{"error": %q}`, err.Error())))
			return
		}

		// Wait before retrying with exponential backoff
		backoff := time.Duration(math.Pow(2, float64(retries))) * time.Second
		time.Sleep(backoff)
	}
}

func processJob(job jobs.Job) error {
	handler, err := Registry.GetHandler(job.Name)
	if err != nil {
		return fmt.Errorf("no handler registered for job type: %s", job.Name)
	}

	// Execute the handler with the job payload
	err = handler(job.Payload)
	if err != nil {
		return fmt.Errorf("job %d failed: %v", job.ID, err)
	}

	return nil
}
