package main

import (
	"context"
	"errors"
	"log"
	"time"

	nursery "github.com/TrevinTeacutter/nursery/pkg/v1"
)

const duration = 2 * time.Second

// This scenario is for things like handling TCP connections that may come and go over the lifetime of a process,
// meaning we don't want `nursery.Wait` to return if there happens to be zero goroutines being managed like you would
// have happen with a `sync.WaitGroup`, we only care about that in the "cleanup" process which happens after the
// context closes.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	defer cancel()

	n := nursery.New(
		nursery.WithContext(ctx),
		nursery.WithWaitForContext(true),
		nursery.WithCloseOnCompletion(false),
		nursery.WithCloseOnError(true),
	)

	n.AddTaskFunc(f, f)

	err := n.Wait()
	if ShouldPanic(err) {
		panic(err)
	}

	log.Println("All jobs done...")
}

func f(ctx context.Context) error {
	for {
		select {
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
		shouldPanic = true
	case errors.Is(err, context.DeadlineExceeded):
		message = "deadline exceeded..."
		shouldPanic = false
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
