package jobs

import (
	"context"

	"github.com/hirakiuc/jobworker-go/models"
)

// BaseJob describe a basic job.
type BaseJob struct {
	message *models.Message
}

// CreateBaseJob create a BasicJob from the message.
func CreateBaseJob(message *models.Message) (Processable, error) {
	return BaseJob{
		message: message,
	}, nil
}

// Process process this BasicJob.
func (j BaseJob) Process(ctx context.Context) error {
	return nil
}
