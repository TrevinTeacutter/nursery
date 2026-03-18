package main

import (
	"context"
	"errors"
	"log"
	"time"

	nursery "github.com/TrevinTeacutter/nursery/pkg/v1"
)

const duration = 2 * time.Second

// This scenario is for long running processes that you want to close together should one exit, for example a worker
// and some health server, if the worker closes regardless of whether an error is present, we likely don't want to keep
// the health server still running. Waiting for context is realistically only necessary if plan to wait when we
// initially have zero goroutines running and add to the nursery after a process waits on the nursery.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	defer cancel()

	n := nursery.New(
		nursery.WithContext(ctx),
		nursery.WithWaitForContext(true),
		nursery.WithCloseOnCompletion(true),
		nursery.WithCloseOnError(true),
	)

	n.AddTaskFunc(longJob, shortJob)

	err := n.Wait()
	if ShouldPanic(err) {
		panic(err)
	}

	log.Println("All jobs done...")
}

func shortJob(ctx context.Context) error {
	timer := time.After(duration / 2)

	for {
		select {
		case <-timer:
			log.Println("job done...")
			return nil
		case <-ctx.Done():
			log.Println("context closed...")
			return nil
		}
	}
}

func longJob(ctx context.Context) error {
	timer := time.After(duration * 2)

	for {
		select {
		case <-timer:
			log.Println("job done...")
			return nil
		case <-ctx.Done():
			log.Println("context closed...")
			return nil
		}
	}
}

func ShouldPanic(err error) bool {
	var (
		message     string
		shouldPanic bool
	)

	switch {
	case errors.Is(err, nursery.CompletionError):
		message = "completed..."
		shouldPanic = false
	case errors.Is(err, context.DeadlineExceeded):
		message = "deadline exceeded..."
		shouldPanic = true
	case errors.Is(err, context.Canceled):
		message = "cancelled..."
		shouldPanic = true
	case errors.Is(err, nursery.ClosedError):
		message = "explicitly closed..."
		shouldPanic = true
	case err == nil:
		message = "exited..."
		shouldPanic = true
	default:
		message = "unexpectedly errored..."
		shouldPanic = true
	}

	log.Println(message)

	return shouldPanic
}
