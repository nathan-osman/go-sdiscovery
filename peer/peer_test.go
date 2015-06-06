package peer

import (
	"testing"
	"time"

	"github.com/nathan-osman/go-sdiscovery/conn"
)

// Test the ping() method.
func Test_peer_ping(t *testing.T) {

	// Create an empty peer.
	times := generateTimes()
	p := &Peer{}

	// Ping the peer twice.
	for i := 0; i < 2; i++ {
		p.ping(&conn.Packet{}, times[0])
	}

	// There should be one address in the peer.
	if len(p.addrs) != 1 {
		t.Fatal("Expected exactly one address")
	}
}

// Test the isExpired() method.
func Test_peer_isExpired(t *testing.T) {

	// Create a peer with an expired address.
	times := generateTimes(2 * time.Second)
	p := &Peer{
		addrs: []*peerAddr{newPeerAddr(nil, times[0])},
	}

	// The peer should have expired.
	if !p.isExpired(time.Second, times[1]) {
		t.Fatal("Peer should be expired")
	}
}
