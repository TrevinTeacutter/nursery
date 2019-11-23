package nursery

import (
	"context"
	"sync"
	"sync/atomic"
)

const (
	CloseError      = basic("cancelled nursery due to closure being requested")
	CompletionError = basic("cancelled nursery due to completed goroutine")
	ClosedError     = basic("nursery is closed")
)

type token struct{}

type Nursery struct {
	// state
	pool      chan token
	counter   atomic.Int64
	context   context.Context
	cancel    context.CancelCauseFunc
	waitGroup sync.WaitGroup
	once      sync.Once

	// configuration
	closeOnFirstCompletion bool
	waitForContext         bool
	limit                  int
}

func New(options ...Option) *Nursery {
	n := &Nursery{
		context: context.Background(),
	}

	for _, option := range options {
		option(n)
	}

	if n.limit > 0 {
		n.pool = make(chan token, n.limit)
	}

	// Given how often it is that these are run with contexts, it would be nice to see this functionality be added into
	// contexts themselves in stdlib. I've attempted to do this but was initially unsuccessful as I wasn't allowing for
	// children/parent relationships when setting them up. Without a context, WaitGroup, ErrGroup are probably more
	// than sufficient.
	// Storing these are less than ideal, but we need a stdlib way of propagating a closure to all goroutines, and this
	// is probably the most accepted way.
	n.context, n.cancel = context.WithCancelCause(n.context)

	return n
}

func (n *Nursery) AddTask(tasks ...Task) {
	for _, task := range tasks {
		select {
		case <-n.context.Done():
			// This prevents adding tasks after this has closed, nurseries are stateful so we need to prevent tasks,
			// getting added after a closure has occurred, even if ideally the task should not do anything because it,
			// should respect the context closure.
			return
		default:
			n.add(task.Run)
		}
	}
}

func (n *Nursery) AddTaskFunc(tasks ...TaskFunc) {
	for _, task := range tasks {
		select {
		case <-n.context.Done():
			// This prevents adding tasks after this has closed, nurseries are stateful so we need to prevent tasks,
			// getting added after a closure has occurred, even if ideally the task should not do anything because it,
			// should respect the context closure.
			return
		default:
			n.add(task)
		}
	}
}

func (n *Nursery) Err() error {
	return context.Cause(n.context)
}

func (n *Nursery) Wait() error {
	select {
	case <-n.context.Done():
		return ClosedError
	default:
	}

	if n.waitForContext {
		<-n.context.Done()
	}

	n.waitGroup.Wait()

	return context.Cause(n.context)
}

func (n *Nursery) Active() int {
	return int(n.counter.Load())
}

func (n *Nursery) Close() {
	n.cancel(CloseError)
}

func (n *Nursery) add(task TaskFunc) {
	if n.pool != nil {
		select {
		case <-n.context.Done():
			return
		case n.pool <- token{}:
			// We never close pool because of this case, select does not respect order of options, so even if
			// context is closed, we are not guaranteed to not try sending a token, once nursery is cleaned up the
			// channel is, so no leaks should occur even if we never close it.
		}
	}

	// It would be cool if `sync.WaitGroup` exposed a method to grab the number of active goroutines it is
	// tracking just so I didn't need a separate atomic value.
	n.waitGroup.Add(1)
	n.counter.Add(1)

	go n.run(task)
}

func (n *Nursery) done() {
	if n.closeOnFirstCompletion {
		defer n.cancel(CompletionError)
	}

	n.waitGroup.Done()
	n.counter.Add(-1)

	if n.pool != nil {
		select {
		case <-n.pool:
		case <-n.context.Done():
		}
	}
}

func (n *Nursery) run(task TaskFunc) {
	defer n.done()

	if err := task(n.context); err != nil {
		// We only care about the first error, while you could potentially aggregate errors, you could potentially end
		// up with several that wrap the context closure error, and in most use cases that I've seen, we only care about
		// the triggering error.
		n.once.Do(func() {
			n.cancel(err)
		})
	}
}
