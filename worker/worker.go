// main/worker/worker.go
package worker

import (
	"fmt"
	"github/Eranmonnie/jqueue/jobs"
	"github/Eranmonnie/jqueue/queue"
	"log"
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

// RegisterHandler adds a job handler function for a specific job type
func (wr *WorkerRegistry) RegisterHandler(jobType string, handler JobHandler) {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	wr.handlers[jobType] = handler
	log.Printf("Registered handler for job type: %s", jobType)

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
func StartWorker(id int) {
	log.Printf("Worker %d started", id)

	for {
		job, err := queue.DequeueJob()
		if err != nil {
			log.Printf("Worker %d error: %v", id, err)
			time.Sleep(2 * time.Second) // Retry delay
			continue
		}

		log.Printf("Worker %d processing job %d of type %s", id, job.ID, job.Name)
		processJob(job)
	}
}

// Process job based on type
func processJob(job jobs.Job) {
	handler, err := Registry.GetHandler(job.Name)
	if err != nil {
		log.Printf("Error processing job %d: %v", job.ID, err)
		// Here you would update the job status to failed
		return
	}

	// Execute the handler with the job payload
	err = handler(job.Payload)
	if err != nil {
		log.Printf("Job %d failed: %v", job.ID, err)
		// Update job status to failed with the error
		return
	}

	log.Printf("Job %d completed successfully", job.ID)
	// Update job status to completed
}
