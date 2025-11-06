# Interview Preparation Guide - Golang

## Operators
* Arithmetic: `+, -, *, /, %, ++, --`
* Comparison: `==, !=, >, <, >=, <=`
* Logical: `&&, ||, !`
* Bitwise: `&, |, ^ (XOR), &^ (clear, AND NOT), <<, >>`
* Assignment: `=, +=, -=, *=, /=, %=, &=, |=, ^=, &^=, <<=, >>=`
* Channel operator: `<-` (send/ receive - based on prefix or suffix of a channel)
* Address: `&` (address of), `*` (dereference)

> It doesn't include pre-increment or pre-decrement operators.
In Go, i++ and i-- are statements, not expressions. This means
they cannot be used within other expressions or function calls
where value is expected. e.g. `fmt.Println(i++)`

## goroutine
It is an independent function launched by go statement, capable of being run concurrently with other goroutines. It is a lightweight thread of execution managed by Go runtime scheduler.
```go
func someFunction() {
    // do something
}
func main() {
    go someFunction()

    go func() {
        fmt.Println("Hello")
    }()
}
```
**Go runtime scheduler:** It is responsible for distributing the runnable goroutines over multiple OD threads. Work sharing and work stealing modes.

The runtime package provides a number of functions that can be used to query and make small changes to the Go routine.
e.g. GOMAXPROCS - to set number of runtime threads that the Go routine will use

## channels

Go channels are fundamental concurrency primitive in Golang,
providing a way for goroutines to communicate and synchronize
their execution. They act as typed conduits through which values
can be sent and received. Channels are designed to be safe for
concurrent use, meaning they handle synchronization automatically,
preventing data races and simplifying concurrent programming.

* **Typed:**
Channels are declared with a specific data type, ensuring that
only values of that type can be sent or received.
* **Blocking by default:**
Sends and receives on unbuffered channels (channels with zero
capacity) are blocking. A send operation will block until a
receiver is ready to receive the value, and a receive operation
will block until a sender sends a value. This blocking behavior
facilitates synchronization between goroutines.
* **Buffered channels:**
Channels can be created with a buffer, allowing them to hold
a certain number of values before blocking. A send to a bufferred
channel will block only if the buffer is full, and a receive will
only block if the buffer is empty.

* **Example - Goroutines synchronization**
```go
package main

import (
	"fmt"
	"time"
)

func worker(done chan bool) {
	fmt.Println("Worker: Starting work...")
	time.Sleep(2 * time.Second) // Simulate work
	fmt.Println("Worker: Work finished.")
	done <- true // Signal that work is done
}

func main() {
	done := make(chan bool)
	go worker(done)
	<-done // Block until the worker signals completion
	fmt.Println("Main: Worker completed.")
}
```

* **Example - Data transfer between Goroutines**
```go
package main

import (
    "fmt"
    "time"
)

func producer(data chan int) {
    for i := 0; i < 5; ++i {
        data <- i
    }
    close(data)
}

func consumer(data chan int) {
    for num := range data {
        fmt.Println("Received: ", num)
    }
}

func main() {
    dataChannel := make(chan int)
    go producer(dataChannel)
    consumer(dataChannel)
}
```

* **Error Handling:**
Channels can be used to communicate errors that occur within
a goroutine back to the main goroutine or other error-handling
routines.

* **Implementing concurrency patterns:**
Channels are integral to implementing various concurrency
patterns like worker pools, fan-out/fan-in, and timeout using
select statements.

## select
Golang's select lets you wait on multiple channel operations.
Combining goroutines and channels with select is powerful
feature of Go.

Example of two channels receiving values.
We can use `select` to await both of these channel values
simultaneously.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(2 * time.Second)
        c1 <- "one"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
    }()

    for range 2 {
        select {
        case msg1 := <-c1:
            fmt.Println("From channel c1:", msg1)
        case msg2 := <-c2:
            fmt.Println("From channel c2:", msg2)
        }
    }
}
```

## How does Go scheduler multiplex goroutines onto OS threads?
The Go sheduler uses `M:N` threading model to multiplex goroutines
onto OS threads, meaning M goroutines are mapped to N OS threads,
where M is typically much larger than N. This allows Golang to run
a large number of goroutines concurrently while maintaining a
relatively small number of OS threads, optimizing resource usage.

Proessors: These are logical CPUs managed by the Go scheduler.
Each processor is associated with an OS thread and maintains a local
run queue of goroutines.

OS Threads: These are kernel-level threads provided by the operating
system. The Go scheduler multiplexes goroutines onto these threads.

### Scheduling
* When a goroutine is created, it is added to a processor's local
run queue.
* The processor selects a goroutine from its queue and assigns it
to its associated OS thread for execution.
* If a processor's local queue is empty, it can "steal" goroutines
from other processor's queues to ensure all threads are kept busy.
* Preemption: The Go scheduler can also preempt(interrupt)
long-running goroutines to ensure fairness and prevent any single
goroutine from monopolizing a thread.

> **Efficiency**: By multiplexing many goroutines onto a smaller
number of threads, Go avoids the overhead of creating and managing
a large number of OS threads, making it efficient for concurrent
operations.

> **Scalability:** The ability to steal work from other processors
ensures that the available threads are kept busy, even when some
proessors have idle goroutines.

> **Concurrency, not just parallelism:** Go's scheduler enables
concurrency (the ability to handle multiple tasks at the same
time), and parallelism (actually running tasks simultaneously
on multiple cores) is achieved by using multiple OS threads.

## Difference between sync.Mutex and channel-based synchronization

In Go, `sync.Mutex` and channel-based synchronization are both
used for managing concurrent access to shared resources, but they
differ in their approach and suitability for various schenarios.
Mutexes provide mutual exclusion, allowing only one goroutine
to access a critical section at a time, while channels facilitate
communication and synchronization between goroutines by passing
values. Generally, channels are preferred for more complex
synchronization patterns and data flow control, while mutexes
are better suited for simple resource protection.

### Mechanism
* `sync.Mutex`: A locking mechanism that provides exclusive access
to a resource. A goroutine acquires the lock (`Lock()`) before
accessing the resource and releases it (`Unlock()`) afterward.
Other goroutines attempting to acquire the lock will be blocked
until it is released.
* **Channel-based synchronization:** Uses channels for communication
and coordination between goroutines. Values are sent (`<-`) and
received (`<-`) on channels to signal events, synchronize actions,
or pass data, allowing for more complex patterns like pipelines.

### Use Cases
* `sync.Mutex`: Ideal for protecting shared variables or data
structures where mutual exclusion is the primary requirement.
Good for situations where performance is critical and the
overhead of channel communication is undesirable.

* Channel-based synchronization: Suited for more complex
scenarios like producer-consumer patterns, pipelines, and
managing concurrent operations with specific coordination requirements.

### Complexity
* `sync.Mutex`: Simpler to understand and use for basic
synchronization needs.
* Channel-based synchronization: Potentially more complex to
implement for advanced patterns but can offer more flexibility
and expressiveness for certain synchronization problems.

## Mutex example
```go
import "sync"

var (
    count int
    mu sync.Mutex
)

func increment() {
    mu.Lock()
    defer mu.Unlock()
    count++
}
```

## Channel example
```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        // Process the job
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= 5; a++ {
        <-results
    }
}
```
## How to avoid goroutine leaks?
A goroutine leak occurs in Go when a goroutine is started but never
properly terminates, remaining active indefinitely and consuming
system resources. These "stuck" goroutines can accumulate over time,
leading to increased memory usage, degraded application performance,
and potentially system crashes.
* **blocking operations:** A goroutine waiting indefinitely on a
channel, mutex, or other synchronization primitive that is never
released or closed.
* **infinite loops:** A goroutine stuck in a loop without a proper
exit condition.
* **Forgotten goroutines:** A goroutine launched without a mechanism
to signal its completion or cancellation, even if its task is no
longer needed.

### To avoid goroutine leaks:
* Use context.Context for cancellation.
Pass a context.Context to goroutines that perform cancellable
operations. Use context.WithCancel or content.WithTimeout to
create contexts that can signal goroutines to stop.
```go
func worker(ctx context.Context, data chan int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker received cancellation signal.")
            return
        case val := <-data:
            fmt.Printf("Worker received: %d\n", val)
        }
    }
}
```

* Ensure channels are properly closed: If a goroutine is waiting
on a channel, ensure that the channel is eventually closed or that
the goroutine has an alternative exit path.
```go
func producer(ch chan int) {
    defer close(ch) // Important to close the channel when done
    for i := 0; i < 5; i++ {
        ch <- i
    }
}
```

* Handle panics gracefully: Use `defer` and `recover` within
goroutines to catch panics and prevent them from causing the
entire program to crash or leaving goroutine in an unhandled state.
```go
func safeTask() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    // ... potentially panicking code ...
}
```
* Ensure that any loops within goroutines have well-defined
conditions for termination.

* Monitor goroutine count: Regularly check `runtime.NumGoroutine()`
to identify potential leaks early. An unexpected increase in the
count can indicate a leak.

* Use `sync.WaitGroup` for tracking completion: If a group of
goroutines needs to complete before proceeding, use `sync.WaitGroup`
to ensure all goroutines finish.
```go
func main() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Goroutine %d completed\n", id)
        }(i)
    }
    wg.Wait() // Wait for all goroutines to finish
    fmt.Println("All goroutines completed.")
}
```

## Use of context.Context for cancellation.
The `context.Context` type in Go serves as a mechanism to manage
deadlines, cancellation signals, and request-scoped values across
API boundaries and between processes. Its primary use for cancellation
involves propagating a signal that indicates an operation should be
stopped or aborted.

* `context.WithCancel(parent Context)`: returns a new `Context` and a
`CancelFunc`. Calling the `CancelFunc` explicitly cancels this `Context`
and all its derived children.
* `context.WithTimeout(parent Context, timeout time.Duration)`: returns
a new `Context` and `CancelFunc`. This context automatically cancels after
the specified timeout duration, or when the `CancelFunc` is called.
* `context.WithDeadline(parent Context, deadline time.Time)`: similar
to `WithTimeout`, but cancels at a specific `deadline` time.

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func longRunningOperation(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Operation completed successfully!")
		return nil
	case <-ctx.Done():
		fmt.Println("Operation cancelled!")
		return ctx.Err() // Return the cancellation error
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure cancel is called to release resources

	err := longRunningOperation(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
```
