package jobs

import "context"

// Processable describe the interface which can be processed by worker.
type Processable interface {
	Process(ctx context.Context) error
}
