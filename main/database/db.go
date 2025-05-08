package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Options  map[string]string
}

// GormDB implements the Database interface using GORM
type GormDB struct {
	config DatabaseConfig
	db     *gorm.DB
}

// NewGormDB creates a new GormDB instance
func NewGormDB(config DatabaseConfig) (*GormDB, error) {
	return &GormDB{
		config: config,
	}, nil
}

// Connect establishes a connection to the database
func (g *GormDB) Connect() error {
	var dialector gorm.Dialector

	switch g.config.Type {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			g.config.Host, g.config.Port, g.config.User, g.config.Password, g.config.Name)

		if opts := g.config.Options; opts != nil {
			if _, ok := opts["sslmode"]; !ok {
				dsn += " sslmode=disable"
			}
			for k, v := range opts {
				dsn += fmt.Sprintf(" %s=%s", k, v)
			}
		}

		dialector = postgres.Open(dsn)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			g.config.User, g.config.Password, g.config.Host, g.config.Port, g.config.Name)
		if opts := g.config.Options; opts != nil {
			for k, v := range opts {
				dsn += fmt.Sprintf("&%s=%s", k, v)
			}
		}
		dialector = mysql.Open(dsn)
	default:
		return fmt.Errorf("unsupported database type: %s", g.config.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	g.db = db
	return nil
}

func (g *GormDB) Close() error {
	DB, err := g.db.DB()
	if err != nil {
		return err
	}
	return DB.Close()
}

func (g *GormDB) Ping() error {
	DB, err := g.db.DB()
	if err != nil {
		return err
	}
	return DB.Ping()
}

func (g *GormDB) Initialize() error {
	return g.db.AutoMigrate(&Job{})
}

func (g *GormDB) GetDB() *gorm.DB {
	return g.db
}

func (g *GormDB) CreateJob(job *Job) error {
	return g.db.Create(job).Error
}

func (g *GormDB) GetJob(id uint) (*Job, error) {
	var job Job
	if err := g.db.First(&job, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("job with ID %d not found", id)
		}
		return nil, err
	}
	return &job, nil
}

func (g *GormDB) UpdateJob(job *Job) error {
	return g.db.Save(job).Error
}

func (g *GormDB) DeleteJob(id uint) error {
	return g.db.Delete(&Job{}, id).Error
}

// GetNextJob gets the next job to process based on priority and creation time
func (g *GormDB) GetNextJob() (*Job, error) {
	var job Job
	if err := g.db.Where("status = ?", JobStatusPending).
		Order("priority desc, created_at asc").
		First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No pending jobs
		}
		return nil, err
	}
	return &job, nil
}

// ListJobs retrieves jobs based on filters
func (g *GormDB) ListJobs(filter JobFilter, limit, offset int) ([]*Job, error) {
	var jobs []*Job
	query := g.db.Model(&Job{})

	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.OwnerID != nil {
		query = query.Where("owner_id = ?", *filter.OwnerID)
	}
	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}

	if err := query.Order("priority desc, created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// MarkJobStarted updates a job as started
func (g *GormDB) MarkJobStarted(id uint) error {
	now := time.Now()
	return g.db.Model(&Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     JobStatusRunning,
			"started_at": now,
		}).Error
}

// MarkJobCompleted updates a job as completed
func (g *GormDB) MarkJobCompleted(id uint, result datatypes.JSON) error {
	now := time.Now()
	return g.db.Model(&Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       JobStatusCompleted,
			"result":       result,
			"completed_at": now,
		}).Error
}

// MarkJobFailed updates a job as failed and increments retry count
func (g *GormDB) MarkJobFailed(id uint, result datatypes.JSON) error {
	// Begin a transaction
	return g.db.Transaction(func(tx *gorm.DB) error {
		// Get current job
		var job Job
		if err := tx.First(&job, id).Error; err != nil {
			return err
		}

		// Increment retries
		job.Retries++

		// Update job status based on retry count
		if job.Retries >= job.MaxRetries {
			job.Status = JobStatusFailed
		} else {
			job.Status = JobStatusPending
		}

		job.Result = result

		// Save changes
		return tx.Save(&job).Error
	})
}

// CountJobs returns the count of jobs matching the given filter
func (g *GormDB) CountJobs(filter JobFilter) (int64, error) {
	var count int64
	query := g.db.Model(&Job{})

	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.OwnerID != nil {
		query = query.Where("owner_id = ?", *filter.OwnerID)
	}
	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
