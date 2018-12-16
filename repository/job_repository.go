package repository

import (
	"github.com/hirakiuc/cloud-pubsub-worker-go/jobs"
	"github.com/hirakiuc/cloud-pubsub-worker-go/models"
)

// JobRepository describe a repository to use Job models.
type JobRepository struct {
}

// NewJobRepository return a JobRepository instance.
func NewJobRepository() (*JobRepository, error) {
	return &JobRepository{}, nil
}

// CreateJob return a kind of Jobbable instance which depend on message.
func (r *JobRepository) CreateJob(msg *models.Message) (jobs.Processable, error) {
	// TODO Create a kind of Jobbable instance which depend on message.
	return jobs.CreateBaseJob(msg)
}
