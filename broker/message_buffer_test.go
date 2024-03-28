package broker

import (
	"github.com/mochi-mqtt/server/v2/packets"
	"sync"
	"testing"
)

func TestBufferInitialization(t *testing.T) {
	size := 5
	cb := newBuffer(size)
	if len(cb.buffer) != size {
		t.Errorf("Expected buffer size of %d, got %d", size, len(cb.buffer))
	}
	if cb.count != 0 {
		t.Errorf("Expected buffer to be empty, count = %d", cb.count)
	}
}

func TestPushBelowCapacity(t *testing.T) {
	size := 5
	cb := newBuffer(size)

	for i := 0; i < size-1; i++ {
		msg := packets.Packet{Payload: []byte{byte(i)}}
		cb.push(msg)
		if cb.count != i+1 {
			t.Errorf("Expected count %d, got %d", i+1, cb.count)
		}
		if cb.buffer[i].Payload[0] != byte(i) {
			t.Errorf("Expected payload %d, got %d", byte(i), cb.buffer[i].Payload[0])
		}
	}
}

// TestPushAtCapacity checks behavior when pushing packets into a full buffer
func TestPushAtCapacity(t *testing.T) {
	size := 3
	cb := newBuffer(size)

	// Fill the buffer to capacity
	for i := 0; i < size; i++ {
		msg := packets.Packet{Payload: []byte{byte(i)}}
		cb.push(msg)
	}

	// push another packet, which should overwrite the oldest packet
	newMsg := packets.Packet{Payload: []byte{99}}
	cb.push(newMsg)

	if cb.count != size {
		t.Errorf("Expected count to remain %d, got %d", size, cb.count)
	}

	// The first packet should be overwritten, so the start should be the second packet
	expectedStartPayload := byte(1)
	if cb.buffer[cb.start].Payload[0] != expectedStartPayload {
		t.Errorf("Expected start packet payload %d, got %d", expectedStartPayload, cb.buffer[cb.start].Payload[0])
	}
}

// TestGetEmptyBuffer verifies behavior when retrieving from an empty buffer
func TestGetEmptyBuffer(t *testing.T) {
	cb := newBuffer(3)
	messages := cb.get()
	if len(messages) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(messages))
	}
}

// TestGetNonEmptyBuffer checks retrieval from a non-empty buffer
func TestGetNonEmptyBuffer(t *testing.T) {
	size := 3
	cb := newBuffer(size)
	expectedPayloads := []byte{0, 1, 2}

	for _, payload := range expectedPayloads {
		msg := packets.Packet{Payload: []byte{payload}}
		cb.push(msg)
	}

	messages := cb.get()
	if len(messages) != len(expectedPayloads) {
		t.Fatalf("Expected %d messages, got %d", len(expectedPayloads), len(messages))
	}

	for i, msg := range messages {
		if msg.Payload[0] != expectedPayloads[i] {
			t.Errorf("Expected payload %d, got %d", expectedPayloads[i], msg.Payload[0])
		}
	}
}

// TestConcurrentAccess tests the buffer's behavior under concurrent access
func TestConcurrentAccess(t *testing.T) {
	size := 100
	cb := newBuffer(size)
	wg := sync.WaitGroup{}

	// Perform concurrent pushes
	concurrentPushes := 1000
	wg.Add(concurrentPushes)
	for i := 0; i < concurrentPushes; i++ {
		go func(i int) {
			defer wg.Done()
			msg := packets.Packet{Payload: []byte{byte(i % 256)}}
			cb.push(msg)
		}(i)
	}

	wg.Wait()

	if cb.count != size {
		t.Errorf("Expected buffer count to be %d, got %d", size, cb.count)
	}
}
