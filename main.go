package main

import (
	"fmt"
	"sync"
	"path"
	"time"
)

const (
	WS_CLOSED = iota
	WS_IDLE
	WS_BUSY
)

type WorkStack struct {
	mutex *sync.Mutex
	items []string
	cond  *sync.Cond
	numWorkers int
	maxWorkers int
	closed bool
	workerFunc func(*WorkStack)
	wg *sync.WaitGroup
}

func (ws *WorkStack) Close() {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	ws.closed = true
	ws.cond.Broadcast()
}

func (ws *WorkStack) Push(item string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if ws.numWorkers < ws.maxWorkers {
		go ws.workerFunc(ws)
		ws.numWorkers++
		ws.wg.Add(1)
	}

	ws.items = append(ws.items, item)
	ws.cond.Signal()
}

func (ws *WorkStack) Next() (string, int) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	for !ws.closed && len(ws.items) == 0 {
		ws.cond.Wait()
	}

	if len(ws.items) == 0 && ws.closed {
		return "", WS_CLOSED
	}

	item := ws.items[len(ws.items)-1]
	ws.items = ws.items[:len(ws.items)-1]

	if len(ws.items) > 0 {
		return item, WS_BUSY
	} else {
		return item, WS_IDLE
	}
}

func worker(ws *WorkStack) {
	fmt.Println("starting worker")
	for {
		item, state := ws.Next()
		if state == WS_CLOSED {
			break
		}

		fmt.Println(item)
		if len(item) < 20 {
			ws.Push(path.Join(item, "L"))
			ws.Push(path.Join(item, "R"))
		}

		if state == WS_IDLE {
			break
		}
	}
	fmt.Println("ending worker")

	ws.mutex.Lock()
	ws.numWorkers--
	ws.mutex.Unlock()

	ws.wg.Done()
}

func main() {
	m := &sync.Mutex{}
	ws := &WorkStack{
		mutex: m,
		cond: sync.NewCond(m),
		workerFunc: worker,
		maxWorkers: 10,
		wg: &sync.WaitGroup{},
	}

	ws.Push("foo")
	time.Sleep(time.Second)
	ws.Push("bar")
	time.Sleep(time.Second)
	ws.Push("baz")

	ws.Close()
	ws.wg.Wait()
}
