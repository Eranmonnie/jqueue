// main/queue/redis.go
package queue

import (
	"context"
	"fmt"
	"github/Eranmonnie/jqueue/config"
	"github/Eranmonnie/jqueue/jobs"
	"log"
	"time"
)

const queueName = "jobqueue"

// Enqueue a job
func EnqueueJob(job jobs.Job) error {
	data, err := job.ToJSON()
	if err != nil {
		return err
	}

	_, err = config.RedisClient.RPush(context.Background(), queueName, data).Result()
	if err != nil {
		return fmt.Errorf("failed to enqueue job: %v", err)
	}

	log.Printf("Enqueued job %d", job.ID)
	return nil
}

func DequeueJob() (jobs.Job, error) {
	result, err := config.RedisClient.BLPop(context.Background(), 0*time.Second, queueName).Result()
	if err != nil {
		return jobs.Job{}, fmt.Errorf("failed to dequeue job: %v", err)
	}

	jobData := result[1]
	job, err := jobs.FromJSON([]byte(jobData))
	if err != nil {
		return jobs.Job{}, err
	}

	return job, nil
}
