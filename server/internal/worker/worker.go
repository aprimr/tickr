package worker

import (
	"context"
	"log/slog"
	"sync"
)

// Job represents task that can be processed in the background
type Job func(ctx context.Context) error

// Pool represents the collection of Jobs that needs to be processed in the background
type Pool struct {
	jobQueue chan Job
	wg       sync.WaitGroup
	logger   *slog.Logger
}

// NewPool returns the instance of Pool with the buffered channel of size queueSize
func NewPool(queueSize int, logger *slog.Logger) *Pool {
	return &Pool{
		jobQueue: make(chan Job, queueSize),
		logger:   logger,
	}
}

// Start launches a number of worker goroutines
func (p *Pool) Start(ctx context.Context, workerCount int) {
	for i := 1; i <= workerCount; i++ {
		p.wg.Add(1)

		// Starts a new go routine to handle a job
		go p.WorkerLoop(ctx, i)
	}
}

// WorkerLoop checks for the jobQueue channel for any jobs and process the job if any
func (p *Pool) WorkerLoop(ctx context.Context, id int) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		// check if any jobs added in the jobQueue channel
		case job, ok := <-p.jobQueue:
			// return if the channel is closed
			if !ok {
				return
			}

			// process the job
			err := job(ctx)
			if err != nil {
				p.logger.Error("background job failed", "worker id", id, "error", err)
			}
		}
	}
}

// Enqueue sends a job to the jobQueue channel
func (p *Pool) Enqueue(job Job) {
	p.jobQueue <- job
}

// Stop closes the jobQueue
// and waits for all the pending worker to finish the current job
func (p *Pool) Stop() {
	close(p.jobQueue)
	p.wg.Wait()
}
