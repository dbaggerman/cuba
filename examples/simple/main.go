package main

import (
	"fmt"
	"path"
	"time"

	"github.com/dbaggerman/cuba"
)

func worker(item string) []string {
	fmt.Println(item)
	subs := []string{}
	if len(item) < 20 {
		subs = []string{
			path.Join(item, "L"),
			path.Join(item, "R"),
		}
	}
	return subs
}

func main() {
	ws := cuba.NewStack(worker)

	ws.Push("foo")
	time.Sleep(time.Second)
	ws.Push("bar")
	time.Sleep(time.Second)
	ws.Push("baz")

	ws.Run()
}
