package main

import (
	"context"
	"log"
	"time"

	nursery "github.com/TrevinTeacutter/nursery/pkg/v1"
)

func main() {
	n := nursery.New(
		nursery.WithWaitForContext(false),
		nursery.WithWaitForCompletion(false),
	)

	n.AddTaskFunc(func(_ context.Context) error {
		time.Sleep(time.Millisecond * 10)
		log.Println("Job 1 done...")

		return nil
	}, func(_ context.Context) error {
		time.Sleep(time.Millisecond * 5)
		log.Println("Job 2 done...")

		return nil
	})

	err := n.Wait()
	if err != nil {
		panic(err)
	}

	log.Println("All jobs done...")
}
