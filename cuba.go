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

// Constructs a new Cuba thread pool.
//
// The worker callback will be called by multiple goroutines in parallel, so is
// expected to be thread safe.
//
// Bucket affects the order that items will be processed in. cuba.NewQueue()
// provides FIFO ordering, while cuba.NewStack() provides LIFO ordered work.
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

// Sets the maximum number of worker goroutines.
//
// Default: runtime.NumCPU() (i.e. the number of CPU cores available)
func (cuba *Cuba) SetMaxWorkers(n int32) {
	cuba.maxWorkers = n
}

// Push an item into the worker pool. This will be scheduled to run on a worker
// immediately.
func (cuba *Cuba) Push(item interface{}) {
	cuba.mutex.Lock()
	defer cuba.mutex.Unlock()

	// The ideal might be to have a fixed pool of worker goroutines which all
	// close down when the work is done.
	// However, since the bucket can drain down to 0 and appear done before the
	// final worker queues more items it's a little complicated.
	// Having a floating pool means we can restart workers as we discover more
	// work to be done, which solves this problem at the cost of a little
	// inefficiency.
	if cuba.numWorkers < cuba.maxWorkers {
		cuba.wg.Add(1)
		go cuba.runWorker()
	}

	cuba.bucket.Push(item)
	cuba.cond.Signal()
}

// Push multiple items into the worker pool.
// 
// Compared to Push() this only aquires the lock once, so may reduce lock
// contention.
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

// Calling Finish() waits for all work to complete, and allows goroutines to shut
// down.
func (cuba *Cuba) Finish() {
	cuba.mutex.Lock()

	cuba.closed = true
	cuba.cond.Broadcast()

	cuba.mutex.Unlock()
	cuba.wg.Wait()
}

func (cuba *Cuba) next() (interface{}, bool) {
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
		item, ok := cuba.next()
		if !ok {
			break
		}

		newItems := cuba.workerFunc(item)
		cuba.PushAll(newItems)
	}
	atomic.AddInt32(&cuba.numWorkers, -1)

	cuba.wg.Done()
}
