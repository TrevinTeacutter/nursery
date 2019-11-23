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
		nursery.WithWaitForContext(true),
		nursery.WithWaitForCompletion(true),
	)
	f := func(ctx context.Context) error {
		timer := time.After(duration)

		for {
			select {
			case <-timer:
				return nil
			case <-ctx.Done():
				return nil
			}
		}
	}
	n.AddTaskFunc(f, f)

	err := n.Wait()
	if ShouldPanic(err) {
		panic(err)
	}

	log.Println("All jobs done...")
}

func ShouldPanic(err error) bool {
	switch {
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, nursery.CompletionError):
		log.Println("Deadline exceeded...")

		return false
	case err == nil,
		errors.Is(err, context.Canceled),
		errors.Is(err, nursery.ClosedError):
		return false
	default:
		return true
	}
}
