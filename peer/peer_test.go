package peer

import (
	"testing"
	"time"

	"github.com/nathan-osman/go-sdiscovery/conn"
)

// Ensure that pings result in the addition of new peer addresses.
func Test_Peer_Ping(t *testing.T) {

	// Create an empty peer.
	times := generateTimes()
	p := &Peer{}

	// Ping the peer twice.
	for i := 0; i < 2; i++ {
		p.Ping(&conn.Packet{}, times[0])
	}

	// There should be one address in the peer.
	if len(p.Addrs) != 1 {
		t.Fatal("Expected exactly one address")
	}
}

// Ensure that the peer expires when the addresses expire.
func Test_Peer_IsExpired(t *testing.T) {

	// Create a peer with an expired address.
	times := generateTimes(2 * time.Second)
	p := &Peer{
		Addrs: []*peerAddr{newPeerAddr(nil, times[0])},
	}

	// The peer should have now expired.
	if !p.IsExpired(time.Second, times[1]) {
		t.Fatal("Peer should be expired")
	}
}
