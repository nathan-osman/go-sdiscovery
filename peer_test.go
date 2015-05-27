package sdiscovery

import (
	"testing"
	"time"
)

// Test that peers are properly expired when no pings are received
func Test_HasExpired(t *testing.T) {

	// Create a new peer with the last ping set to the current time
	peer := &Peer{
		Addresses: []PeerAddress{
			{
				IP:         nil,
				lastPacket: time.Now(),
			},
		},
	}

	// The peer should not have expired
	if peer.HasExpired(time.Second) {
		t.Fatal("Peer expired unexpectedly")
	}

	// Set the last ping to two seconds ago
	peer.Addresses[0].lastPacket = time.Now().Add(-2 * time.Second)

	// The peer should now have expired
	if !peer.HasExpired(time.Second) {
		t.Fatal("Peer has not expired")
	}
}
