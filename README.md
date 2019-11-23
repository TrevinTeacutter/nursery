# nursery: structured concurrency in Go

## Installation
```bash
go get -u github.com/TrevinTeacutter/nursery
```

[Notes on structured concurrency, or: Go statement considered harmful](https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/#nurseries-a-structured-replacement-for-go-statements) is an article that compares the dangers of goto with the go statement.

While I don't necessarily agree with the entire content, I can understand the value in having more guarantees and structure on when using goroutines but also potentially for observability reason (for example knowning where in your code you have the most use of goroutines or anything of that nature).

While Go does have some places in the stdlib to help with this such as `sync.WaitGroup` or `golang.org/x/sync.ErrGroup` it can be nice to have better syntactic sugar around these concepts.

For the purpose of this package, there's really two levers of control around when and how a nursery is closed:
* Waiting for closure of a context
* Race for any task to complete

Racing for any task to complete is often a pattern for if you have a bunch of long running goroutines that all should close out if one exits early but does so without error. Say for example an HTTP server offering observability while a job runs, you probably want both decoupled but running one without the other may not make the most sense.

Waiting for closure of a context is useful for bounded jobs but with an insurance policy that there will be a clean exit with no work left hanging. An example of this would be a TCP listener that spawns consumers for each incoming connection. Having 0 jobs is perfectly acceptable so `sync.WaitGroup` semantics on their own aren't good enough.

These could be use together for long running processes as well if you leverage context for closure anywhere in your code.

These can also be avoided entirely if you know some logic is short lived but you need to do the logic concurrently, you will only return once there is an error or all work completed successfully.

Keep in mind though there is no way to hijack or force-stop a goroutine in the go runtime so if someone does not respect context closure this will not buy you the safety you may want so beware.
