package jobs

import (
	"encoding/json"
	"github/Eranmonnie/jqueue/database"
)

type Job struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Status     database.Status `json:"status"`
	Priority   int             `json:"priority"`
	Payload    interface{}     `json:"payload"` // e.g. email content, file path"`
	Retries    int             `json:"retries"`
	MaxRetries int             `json:"max_retries"`
}

func NewJob(id uint, Name string, priority int, MaxRetries int, Payload interface{}) Job {
	job := Job{
		ID:         id, // ID will be set by the database or us later
		Name:       Name,
		Status:     database.JobStatusPending,
		Priority:   priority,
		Payload:    Payload,
		Retries:    0,
		MaxRetries: MaxRetries,
	}
	return job
}

func (j Job) ToJSON() ([]byte, error) {
	return json.Marshal(j)
}

// Deserialize JSON to Job
func FromJSON(data []byte) (Job, error) {
	var job Job
	err := json.Unmarshal(data, &job)
	return job, err
}
