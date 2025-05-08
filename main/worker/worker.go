// main/worker/worker.go
package worker

import (
	"github/Eranmonnie/jqueue/main/jobs"
	"github/Eranmonnie/jqueue/main/queue"

	"log"
	"time"
)

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
	switch job.Name {
	case "email":
		log.Printf("Sending email with payload: %s", job.Payload)
	case "compress":
		log.Printf("Compressing file with payload: %s", job.Payload)
	default:
		log.Printf("Unknown job type: %s", job.Name)
	}
}
