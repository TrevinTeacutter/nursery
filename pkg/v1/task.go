package nursery

import "context"

type Task interface {
	Run(ctx context.Context) error
}

type TaskFunc func(ctx context.Context) error
