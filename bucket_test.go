package cuba

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStackBucket(t *testing.T) {
	stack := NewStack()
	assert.True(t, stack.Empty(), "Newly constructed stack should be empty")

	stack.Push("abc")
	stack.Push("def")
	stack.Push("ghi")
	assert.False(t, stack.Empty(), "Stack should not be empty after push")

	var item string

	item = stack.Pop().(string)
	assert.Equal(t, "ghi", item, "Stack should return items in LIFO order")
	assert.False(t, stack.Empty(), "Stack should not be empty before popping all items")

	item = stack.Pop().(string)
	assert.Equal(t, "def", item, "Stack should return items in LIFO order")
	assert.False(t, stack.Empty(), "Stack should not be empty before popping all items")

	item = stack.Pop().(string)
	assert.Equal(t, "abc", item, "Stack should return items in LIFO order")
	assert.True(t, stack.Empty(), "Stack should be empty after popping all items")
}

func TestQueueBucket(t *testing.T) {
	queue := NewQueue()
	assert.True(t, queue.Empty(), "Newly constructed queue should be empty")

	queue.Push("abc")
	queue.Push("def")
	queue.Push("ghi")
	assert.False(t, queue.Empty(), "Queue should not be empty after push")

	var item string

	item = queue.Pop().(string)
	assert.Equal(t, "abc", item, "Queue should return items in FIFO order")
	assert.False(t, queue.Empty(), "Queue should not be empty before popping all items")

	item = queue.Pop().(string)
	assert.Equal(t, "def", item, "Queue should return items in FIFO order")
	assert.False(t, queue.Empty(), "Queue should not be empty before popping all items")

	item = queue.Pop().(string)
	assert.Equal(t, "ghi", item, "Queue should return items in FIFO order")
	assert.True(t, queue.Empty(), "Queue should be empty after popping all items")
}
