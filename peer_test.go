package sdiscovery

import (
	"testing"
	"time"
)

// Test that addresses are properly expired when no pings are received
func Test_Update(t *testing.T) {

	// Create a new peer with the last ping set to the current time
	peer := &Peer{
		Addresses: []PeerAddress{
			{
				IP:         nil,
				lastPacket: time.Now(),
			},
		},
	}

	// Update the peer and ensure the address has not expired
	peer.Update(time.Second)
	if len(peer.Addresses) != 1 {
		t.Fatal("Address expired unexpectedly")
	}

	// Change the time of last packet to two seconds ago
	peer.Addresses[0].lastPacket = time.Now().Add(-2 * time.Second)

	// Update the peer again and ensure the address expired
	peer.Update(time.Second)
	if len(peer.Addresses) != 0 {
		t.Fatal("Address has not expired")
	}
}

// Test that peers are properly expired when no addresses remain
func Test_HasExpired(t *testing.T) {

	// Create a new peer with a single address
	peer := &Peer{
		Addresses: []PeerAddress{{}},
	}

	// The peer should not have expired
	if peer.HasExpired() {
		t.Fatal("Peer expired unexpectedly")
	}

	// Remove the address
	peer.Addresses = []PeerAddress{}

	// The peer should now have expired
	if !peer.HasExpired() {
		t.Fatal("Peer has not expired")
	}
}
