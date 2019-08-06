package cuba

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type CubaFunc func(interface{}) []interface{}

type Cuba struct {
	mutex      *sync.Mutex
	bucket     Bucket
	cond       *sync.Cond
	numWorkers int32
	maxWorkers int32
	closed     bool
	workerFunc CubaFunc
	wg         *sync.WaitGroup
}

func New(worker CubaFunc, bucket Bucket) *Cuba {
	m := &sync.Mutex{}
	return &Cuba{
		mutex:      m,
		bucket:     bucket,
		cond:       sync.NewCond(m),
		workerFunc: worker,
		maxWorkers: int32(runtime.NumCPU()),
		wg:         &sync.WaitGroup{},
	}
}

func (cuba *Cuba) Close() {
	cuba.mutex.Lock()
	defer cuba.mutex.Unlock()

	cuba.closed = true
	cuba.cond.Broadcast()
}

func (cuba *Cuba) Run() {
	cuba.Close()
	cuba.wg.Wait()
}

func (cuba *Cuba) Push(item interface{}) {
	cuba.mutex.Lock()
	defer cuba.mutex.Unlock()

	if cuba.numWorkers < cuba.maxWorkers {
		cuba.wg.Add(1)
		go cuba.runWorker()
	}

	cuba.bucket.Push(item)
	cuba.cond.Signal()
}

func (cuba *Cuba) PushAll(items []interface{}) {
	cuba.mutex.Lock()
	defer cuba.mutex.Unlock()

	for i := 0; i < len(items); i++ {
		if cuba.numWorkers >= cuba.maxWorkers {
			break
		}
		cuba.wg.Add(1)
		go cuba.runWorker()
	}

	cuba.bucket.PushAll(items)
	cuba.cond.Broadcast()
}

func (cuba *Cuba) Next() (interface{}, bool) {
	cuba.mutex.Lock()
	defer cuba.mutex.Unlock()

	for cuba.bucket.Empty() {
		if cuba.closed {
			return nil, false
		}
		cuba.cond.Wait()
	}

	item := cuba.bucket.Pop()

	return item, true
}

func (cuba *Cuba) runWorker() {
	atomic.AddInt32(&cuba.numWorkers, 1)
	for {
		item, ok := cuba.Next()
		if !ok {
			break
		}

		newItems := cuba.workerFunc(item)
		cuba.PushAll(newItems)
	}
	atomic.AddInt32(&cuba.numWorkers, -1)

	cuba.wg.Done()
}
