package shared

import (
	"testing"
	"time"
)

func TestNewBroadcaster(t *testing.T) {
	b := NewBroadcaster()
	if b == nil {
		t.Fatal("NewBroadcaster returned nil")
	}
	// Check that the consumers map is initialized and empty.
	if b.consumers == nil {
		t.Fatal("Expected consumers map to be initialized")
	}
	if len(b.consumers) != 0 {
		t.Fatal("Expected consumers map to be empty")
	}
	// Check that the input channel is initialized.
	if b.input == nil {
		t.Fatal("Expected input channel to be initialized")
	}
}

func TestBroadcasting(t *testing.T) {
	b := NewBroadcaster()
	consumer := b.Register()
	defer b.Unregister(consumer)

	// Send a message using Send() which sends "update".
	go b.Send()

	select {
	case msg := <-consumer:
		if msg != "update" {
			t.Fatalf("Expected 'update', got %s", msg)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Did not receive broadcast message within timeout")
	}
}
