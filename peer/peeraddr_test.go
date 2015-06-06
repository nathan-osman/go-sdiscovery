package peer

import (
	"container/ring"
	"net"
	"testing"
	"time"
)

var (
	testIP1 = net.IPv4(192, 168, 1, 1)
	testIP2 = net.IPv4(192, 168, 1, 2)

	testTime1 = time.Now()
	testTime2 = testTime1.Add(time.Second)
)

// Count the number of valid elements in a ring.
func validElementsInRing(r *ring.Ring) int {
	i := 0
	r.Do(func(element interface{}) {
		if _, ok := element.(time.Time); ok {
			i++
		}
	})
	return i
}

// Ensure that the ping() method results in the correct behavior.
func Test_peerAddr_ping(t *testing.T) {

	// Create a new peerAddr and confirm that it contains one item.
	p := newPeerAddr(nil, testTime1)
	if validElementsInRing(p.lastPing) != 1 {
		t.Fatal("Expected one element in ring")
	}

	// Register a ping and confirm that the ring now contains two items.
	p.ping(testTime1)
	if validElementsInRing(p.lastPing) != 2 {
		t.Fatal("Expected two elements in ring")
	}
}

// Ensure that duration reports the correct difference between the first and
// last packet in the ring.
func Test_peerAddr_duration(t *testing.T) {

	// Create a peerAddr.
	p := newPeerAddr(nil, testTime1)

	// Ping the address five more times.
	for i := 0; i < 5; i++ {
		p.ping(testTime2)
	}

	// The duration should be one hour.
	if p.duration() != time.Second {
		t.Fatal("Duration does not match")
	}
}

// Ensure that the address expires when the duration is exceeded.
func Test_peerAddr_isExpired(t *testing.T) {

	// Create a new peerAddr with the first time.
	p := newPeerAddr(nil, testTime1)

	// Assuming a timeout of one second, the address should have expired.
	if !p.isExpired(500*time.Millisecond, testTime2) {
		t.Fatal("Address should have expired")
	}
}
