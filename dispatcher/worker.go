package dispatcher

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hirakiuc/jobworker-go/common"
	"go.uber.org/zap"
)

// Worker describe a job processor.
type Worker struct {
	ID uuid.UUID

	// reqChan is NOT owned by Worker, this is used to receive job.
	reqChan chan JobRequest

	// terminateChan can be closed in order to cleanly shutdown this worker.
	terminateChan chan struct{}

	// closedChan is closed by the run() goroutine when it exists.
	closedChan chan struct{}
}

// NewWorker return the worker instance for the Job.
func NewWorker(reqChan chan JobRequest) *Worker {
	return &Worker{
		ID:            uuid.New(),
		reqChan:       reqChan,
		terminateChan: make(chan struct{}),
		closedChan:    make(chan struct{}),
	}
}

// Start starts the worker process.
func (w *Worker) Start() {
	go func() {
		err := w.run()
		if err != nil {
			logger := w.getLogger()
			logger.Error("Worker Failed",
				zap.Error(err),
			)
		}
	}()
}

func (w *Worker) getLogger() *zap.Logger {
	logger := common.GetLogger()
	return logger.With(
		zap.String("WorkerID", w.ID.String()),
	)
}

func (w *Worker) run() error {
	logger := w.getLogger()

	for {
		logger.Info("Worker Loop Start")
		w.HookPrepare()

		logger.Info("Listening next job...")
		select {
		case job, ok := <-w.reqChan:
			if !ok {
				w.HookCleanup()
				return errors.New("reqChan is closed")
			}

			err := w.Process(job)
			if err != nil {
				logger.Error("Job failed", zap.Error(err))
			}

		case _, ok := <-w.terminateChan:
			logger.Info("Worker termination detected")
			if !ok {
				w.HookCleanup()
				return errors.New("w.terminateChan is closed")
			}

			w.HookCleanup()
			return nil
		}

		// TODO sleep here
	}
}

// Process process the job with the context.
func (w *Worker) Process(job JobRequest) error {
	ctx := context.Background()
	return job.Process(ctx)
}

// HookPrepare make the worker ready
func (w *Worker) HookPrepare() {
	logger := common.GetLogger()
	logger.Info("Worker HookPrepare start")
	logger.Info("Worker HookPrepare end")
}

// HookCleanup clean up the worker resources.
func (w *Worker) HookCleanup() {
	// TODO: cleanup tasks

	w.closedChan <- struct{}{}
}

// Terminate ...
func (w *Worker) Terminate() {
	logger := common.GetLogger()
	logger.Info("Worker Terminate Start")

	select {
	case w.terminateChan <- struct{}{}:
		logger.Info("Worker Terminate Send terminate event")
	default:
		logger.Info("Worker Terminate NoCapacity")
	}

	// TODO cleanup channels and resources
	logger.Info("Worker Terminate End")
}

// WaitUntilTerminated wait until this worker terminated.
func (w *Worker) WaitUntilTerminated() {
	logger := common.GetLogger()
	logger.Info("Worker WaitUntilTerminated called")

	_, ok := <-w.closedChan
	if !ok {
		// TODO: log the channel already closed.
	}
}
