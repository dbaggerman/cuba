package main

import (
	"os"
	"fmt"
	"path"
	"time"

	"github.com/dbaggerman/cuba"
)

func worker(itemIf interface{}) []interface{} {
	item := itemIf.(string)
	fmt.Fprintf(os.Stderr, "Item: %s\n", item)
	subs := []interface{}{}
	if len(item) < 20 {
		subs = []interface{}{
			path.Join(item, "L"),
			path.Join(item, "R"),
		}
	}
	return subs
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
