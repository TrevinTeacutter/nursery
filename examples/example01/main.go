package main

import (
	"context"
	"log"
	"time"

	nursery "github.com/TrevinTeacutter/nursery/pkg/v1"
)

const duration = 2 * time.Second

// This scenario is the more traditional use case of a `sync.WaitGroup` where you want to do some work concurrently,
// before joining back up to continue along. We want to give all workers time to do things unless they error out, thus
// why we do not trigger closure on completion. Since this is meant to be temporary, we also do not wait for the
// context to close.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	defer cancel()

	n := nursery.New(
		nursery.WithContext(ctx),
		nursery.WithWaitForContext(false),
		nursery.WithCloseOnCompletion(false),
		nursery.WithCloseOnError(true),
	)

	//n.AddTaskFunc(job1, job2)

	err := n.Wait()
	if err != nil {
		panic(err)
	}

	log.Println("All jobs done...")
}

func job1(_ context.Context) error {
	time.Sleep(time.Millisecond * 10)
	log.Println("Job 1 done...")

	return nil
}

func job2(_ context.Context) error {
	time.Sleep(time.Millisecond * 5)
	log.Println("Job 2 done...")

	return nil
}
