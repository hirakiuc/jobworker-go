package dispatcher

import (
	"os"
	"os/signal"
	"time"

	"github.com/hirakiuc/jobworker-go/client"
	"github.com/hirakiuc/jobworker-go/common"
	"github.com/hirakiuc/jobworker-go/repository"
	"go.uber.org/zap"
)

// Config describe the dispatcher config.
type Config struct {
	NumOfWorkers int

	// RateLimit
	JobRate  time.Duration // 1 job/sec = time.Second/1
	JobBurst int
}

// Dispatcher describe the job dispatcher.
type Dispatcher struct {
	config *Config

	pubSubClient  *client.PubSubClient
	jobRepository *repository.JobRepository
	pool          *Pool
	rateLimiter   *RateLimit

	interrupted bool
}

// NewDispatcher return a new Dispatcher instance.
func NewDispatcher(config *Config) (*Dispatcher, error) {
	pubSubClient, err := client.NewPubSubClient()
	if err != nil {
		return nil, err
	}

	jobRepository, err := repository.NewJobRepository()
	if err != nil {
		return nil, err
	}

	pool, err := NewPool(config.NumOfWorkers)
	if err != nil {
		return nil, err
	}

	rateLimiter := NewRateLimit(config.JobRate, config.JobBurst)

	return &Dispatcher{
		config:        config,
		pubSubClient:  pubSubClient,
		jobRepository: jobRepository,
		pool:          pool,
		rateLimiter:   rateLimiter,
		interrupted:   false,
	}, nil
}

// Start begins to dispatching Messages
func (d *Dispatcher) Start() error {
	logger := common.GetLogger()

	logger.Info("Increase workers",
		zap.Int("NumOfWorkers", d.config.NumOfWorkers),
	)
	d.pool.SetSize(d.config.NumOfWorkers)

	// Handle interruption
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		sig := <-c
		logger.Info("Signal received", zap.String("signal", sig.String()))
		close(c)

		// TODO This should be set atomically.
		d.interrupted = true
	}()

	// Start rateLimitter
	d.rateLimiter.Start()
	defer d.rateLimiter.Stop()

	for {
		// RateLimit
		d.rateLimiter.Throttle()

		// TODO: Handle CTRL-C interruption
		if d.interrupted {
			logger.Info("Detect interruption.")
			return d.Stop()
		}

		logger.Info("Start GetMessage")
		msg, err := d.pubSubClient.GetMessage()
		if err != nil {
			logger.Error("Failed to fetch message", zap.Error(err))
			continue
		}
		if msg == nil {
			continue
		}

		logger.Info("Start CreateJob")
		job, err := d.jobRepository.CreateJob(msg)
		if err != nil {
			logger.Error("Failed to create Job", zap.Error(err))
			// TODO: This `msg` should be send to CloudPubSub for retry.
			continue
		}

		logger.Info("Send the job to worker")
		// TODO: Send the job to idle worker.
		err = d.pool.SendJobAsync(job)
		if err != nil {
			logger.Error("Failed to process the job", zap.Error(err))
			// TODO: This `job` should be send to CloudPubSub for retry.
			continue
		}
	}
}

// Stop terminate all of workers gracefully.
func (d *Dispatcher) Stop() error {
	d.pool.SetSize(0)

	// TODO: Cleanup if needed.
	return nil
}
