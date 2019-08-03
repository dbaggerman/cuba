package cuba

import (
	"sync"
	"runtime"
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
	workerFunc func(string) []string
	wg *sync.WaitGroup
}

func NewStack(worker func(string) []string) *WorkStack {
	m := &sync.Mutex{}
	return &WorkStack{
		mutex: m,
		cond: sync.NewCond(m),
		workerFunc: worker,
		maxWorkers: runtime.NumCPU(),
		wg: &sync.WaitGroup{},
	}
}

func (ws *WorkStack) Close() {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	ws.closed = true
	ws.cond.Broadcast()
}

func (ws *WorkStack) Run() {
	ws.Close()
	ws.wg.Wait()
}

func (ws *WorkStack) Push(item string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if ws.numWorkers < ws.maxWorkers {
		go ws.runWorker()
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

func (ws *WorkStack) runWorker() {
	for {
		item, state := ws.Next()
		if state == WS_CLOSED {
			break
		}

		newItems := ws.workerFunc(item)
		for _, newItem := range newItems {
			ws.Push(newItem)
		}

		if state == WS_IDLE {
			break
		}
	}

	ws.mutex.Lock()
	ws.numWorkers--
	ws.mutex.Unlock()

	ws.wg.Done()
}
