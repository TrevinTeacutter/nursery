package main

import (
	"context"
	"errors"
	"log"
	"time"

	nursery "github.com/TrevinTeacutter/nursery/pkg/v1"
)

const duration = 2 * time.Second

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	defer cancel()

	n := nursery.New(
		nursery.WithContext(ctx),
		nursery.WithWaitForContext(false),
		nursery.WithCloseOnCompletion(false),
		nursery.WithCloseOnError(true),
	)

	n.AddTaskFunc(nursery.Recovery(f))

	err := n.Wait()
	if ShouldPanic(err) {
		panic(err)
	}

	log.Println("All jobs done...")
}

func f(_ context.Context) error {
	panic("oops")
}

func ShouldPanic(err error) bool {
	var (
		message     string
		shouldPanic bool
	)

	switch {
	case errors.Is(err, &nursery.ErrPanic{}):
		message = "panic occurred..."
		shouldPanic = false
	case errors.Is(err, nursery.CompletionError):
		message = "completed..."
		shouldPanic = true
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
