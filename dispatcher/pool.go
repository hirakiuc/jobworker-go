package dispatcher

// Pool keeps Workers
// Workers always running.
// Pool -(Job)-> Worker receive via channel
// Worker process the job
//

import (
	"context"
	"sync"

	"github.com/hirakiuc/jobworker-go/jobs"
)

// JobRequest describe a request instance to be processed by worker.
type JobRequest struct {
	job jobs.Processable
}

func (req JobRequest) Process(ctx context.Context) error {
	return req.job.Process(ctx)
}

// Pool describe a pool of workers.
type Pool struct {
	queuedJobs int64

	workers []*Worker

	reqChan chan JobRequest

	workerMux sync.Mutex
}

// NewPool create a Pool instance
func NewPool(size int) (*Pool, error) {
	pool := &Pool{
		queuedJobs: 0,
		reqChan:    make(chan JobRequest),
	}

	pool.SetSize(0)

	return pool, nil
}

// SendJobAsync send the job to idle worker.
func (pool *Pool) SendJobAsync(job jobs.Processable) error {
	pool.reqChan <- JobRequest{
		job: job,
	}

	return nil
}

// SendJobSync process the job with new worker.
func (pool *Pool) SendJobSync(job *jobs.Processable) error {
	// TODO: Create new worker and process the job here.
	return nil
}

// SetSize configure thw number of workers in this Pool.
func (pool *Pool) SetSize(size int) {
	pool.workerMux.Lock()
	defer pool.workerMux.Unlock()

	current := len(pool.workers)
	if current == size {
		return
	}

	if size > current {
		// Add extra workers if size > current
		for i := current; i < size; i++ {
			worker := NewWorker(pool.reqChan)
			worker.Start()
			pool.workers = append(pool.workers, worker)
		}

		return
	}

	// size < current
	// Asynchronously stop all workers (size ..(current))
	for i := size; i < current; i++ {
		pool.workers[i].Terminate()
	}

	// Synchronously wait for the workers who requested to stop.
	for i := 0; i < (current - size); i++ {
		pool.workers[i].WaitUntilTerminated()
	}
}
