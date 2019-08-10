package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/dbaggerman/cuba"
)

func worker(handle *cuba.Handle) {
	item := handle.Item().(string)
	fmt.Fprintf(os.Stderr, "Item: %s\n", item)
	if len(item) < 20 {
		handle.Push(path.Join(item, "L"))
		handle.Push(path.Join(item, "R"))
	}
}

func main() {
	ws := cuba.New(worker, cuba.NewQueue())

	ws.Push("foo")
	time.Sleep(time.Second)

	ws.Push("bar")
	time.Sleep(time.Second)

	ws.Push("baz")

	ws.Finish()
}
