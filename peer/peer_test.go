package peer

import (
	"testing"
	"time"

	"github.com/nathan-osman/go-sdiscovery/comm"
)

// Ensure that pings result in the addition of new peer addresses.
func Test_Peer_Ping(t *testing.T) {

	// Create an empty peer.
	p := &Peer{}

	// Ping the peer twice.
	for i := 0; i < 2; i++ {
		p.Ping(&comm.Packet{}, testTime1)
	}

	// There should be one address in the peer.
	if len(p.addrs) != 1 {
		t.Fatal("Expected exactly one address")
	}
}

// Ensure that Addrs() returns a properly sorted slice of addresses.
func Test_Peer_Addrs(t *testing.T) {

	// Create a peer with two addresses one second apart.
	p := &Peer{
		addrs: peerSlice{
			newPeerAddr(testIP1, testTime1),
			newPeerAddr(testIP2, testTime2),
		},
	}

	// Obtain the list of addresses.
	addrs := p.Addrs()

	// Ensure that two items are present.
	if len(addrs) != 2 {
		t.Fatal("Expected exactly two addresses")
	}

	// Ensure that the order is correct
	if !addrs[0].Equal(testIP1) || !addrs[1].Equal(testIP2) {
		t.Fatal("Addresses sorted incorrectly")
	}
}

// Ensure that the peer expires when the addresses expire.
func Test_Peer_IsExpired(t *testing.T) {

	// Create a peer with an expired address.
	p := &Peer{
		addrs: peerSlice{newPeerAddr(testIP1, testTime1)},
	}

	// The peer should have now expired.
	if !p.IsExpired(500*time.Millisecond, testTime2) {
		t.Fatal("Peer should be expired")
	}
}
