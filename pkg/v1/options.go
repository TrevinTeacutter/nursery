package nursery

import "context"

type Option func(*Nursery)

func WithContext(value context.Context) Option {

	return func(n *Nursery) {
		n.context = value
	}
}

func WithLimit(value int) Option {
	return func(n *Nursery) {
		n.limit = value
	}
}

func WithWaitForContext(value bool) Option {
	return func(n *Nursery) {
		n.waitForContext = value
	}
}

func WithWaitForCompletion(value bool) Option {
	return func(n *Nursery) {
		n.closeOnFirstCompletion = value
	}
}
