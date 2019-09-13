package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dbaggerman/cuba"
)

// Since this always queues up another item, if we call pool.Finish() it will
// wait forever for all the work to complete.
// Using Abort() means that the items pushed by the worker will be dropped
// instead of being added to the pool.
func worker(handle *cuba.Handle) {
	n := handle.Item().(int)
	fmt.Fprintf(os.Stderr, "Item: %d\n", n)
	time.Sleep(100 * time.Millisecond)
	handle.Push(n + 1)
}

func main() {
	ws := cuba.New(worker, cuba.NewQueue())

	ws.Push(0)

	time.Sleep(time.Second)
	ws.Abort()
}
