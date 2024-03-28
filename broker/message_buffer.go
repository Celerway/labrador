package broker

import (
	"github.com/mochi-mqtt/server/v2/packets"
	"sync"
)

// circularBuffer is a fixed-size buffer that holds the last N elements
// in a circular manner.
type circularBuffer struct {
	buffer []packets.Packet // The fixed-size buffer
	size   int              // Maximum number of messages to store
	start  int              // Points to the start of the buffer
	end    int              // Points to the end of the buffer
	count  int              // Current number of elements in the buffer
	mu     sync.Mutex       // Ensures concurrent access to the buffer is safe
}

// newBuffer creates a new circularBuffer with a given size.
func newBuffer(size int) *circularBuffer {
	return &circularBuffer{
		buffer: make([]packets.Packet, size),
		size:   size,
	}
}

// push adds a new message to the buffer, overwriting the oldest message if the buffer is full.
func (cb *circularBuffer) push(msg packets.Packet) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.count == cb.size {
		// Buffer is full; increment start to overwrite the oldest element
		cb.start = (cb.start + 1) % cb.size
	} else {
		// Buffer is not full; increment count
		cb.count++
	}

	// Add the new message and adjust the end pointer
	cb.buffer[cb.end] = msg
	cb.end = (cb.end + 1) % cb.size
}

// get returns a slice of messages currently in the buffer, in the order they were added.
func (cb *circularBuffer) get() []packets.Packet {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	messages := make([]packets.Packet, cb.count)
	for i := 0; i < cb.count; i++ {
		index := (cb.start + i) % cb.size
		messages[i] = cb.buffer[index]
	}
	return messages
}
