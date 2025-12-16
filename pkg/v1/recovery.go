package nursery

import (
	"context"
	"fmt"
	"runtime/debug"
)

var (
	_ error = (*ErrPanic)(nil)
)

// ErrPanic is an `error` used to contain any useful information about a panic that we recovered from.
type ErrPanic struct {
	value      any
	stacktrace []byte
}

func (e *ErrPanic) Error() string {
	return fmt.Sprintf("panic: %v\nstacktrace:\n%s\n", e.value, e.stacktrace)
}

// Is overrides the default behavior for `errors.Is` to ensure that typing is all that matters for `errors.Is`
// comparison as `errors.As` is much less ergonomic for guards based on this.
func (e *ErrPanic) Is(target error) bool {
	switch target.(type) {
	case *ErrPanic:
		return true
	default:
		return false
	}
}

// Recovery wraps a task function with a panic recovery mechanism. On panic, an `ErrPanic` is returned with the value
// panicked with and a stacktrace from the point of recovery.
func Recovery(f TaskFunc) TaskFunc {
	return func(ctx context.Context) (err error) {
		defer func() {
			if value := recover(); value != nil {
				err = &ErrPanic{value, debug.Stack()}
			}
		}()

		return f(ctx)
	}
}
