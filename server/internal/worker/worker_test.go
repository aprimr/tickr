package worker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	pool := NewPool(10, logger)

	ctx := t.Context()

	// Create a pool of 2 workers
	pool.Start(ctx, 2)

	var count int32
	var jobCount int32 = 20

	// Enqueue 20 jobs in the pool
	for i := range jobCount {
		pool.Enqueue(func(ctx context.Context) error {
			fmt.Printf("added job %v\n", i)
			atomic.AddInt32(&count, 1)

			return nil
		})
	}

	// Wait for all jobs to finish
	pool.Stop()

	if count != jobCount {
		t.Errorf("expected %d jobs to be executed, but got %d", jobCount, count)
	}
}
