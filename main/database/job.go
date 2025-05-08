package database

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Status string

// Job Statuses as constants
const (
	JobStatusPending   Status = "pending"
	JobStatusRunning   Status = "running"
	JobStatusCompleted Status = "completed"
	JobStatusFailed    Status = "failed"
	JobStatusCancelled Status = "cancelled"
)

// Job represents a job in the queue
type Job struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Status      Status         `gorm:"default:pending" json:"Status"`
	Priority    int            `gorm:"default:0" json:"priority"`
	Payload     datatypes.JSON `json:"payload"`
	Result      datatypes.JSON `json:"result"`
	Retries     int            `gorm:"default:0" json:"retries"`
	MaxRetries  int            `gorm:"default:3" json:"max_retries"`
	StartedAt   *time.Time     `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
}

// JobFilter represents filters for querying jobs
type JobFilter struct {
	ID       *uint   `json:"id,omitempty"`
	Status   *Status `json:"Status,omitempty"`
	Priority *int    `json:"priority,omitempty"`
	OwnerID  *string `json:"owner_id,omitempty"`
	Name     *string `json:"name,omitempty"`
}
