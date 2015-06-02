package sdiscovery

import (
	"container/ring"
	"testing"
	"time"
)

// Count the number of valid elements in a ring
func validElementsInRing(r *ring.Ring) int {
	i := 0
	r.Do(func(element interface{}) {
		if _, ok := element.(time.Time); ok {
			i++
		}
	})
	return i
}

// Test the ping() method
func Test_ping(t *testing.T) {

	// Create a new peerAddr and confirm that it contains one item
	p := newPeerAddr(nil, time.Now())
	if validElementsInRing(p.lastPing) != 1 {
		t.Fatal("Expected one element in ring")
	}

	// Register a ping and confirm that the ring now contains two items
	p.ping(time.Now())
	if validElementsInRing(p.lastPing) != 2 {
		t.Fatal("Expected two elements in ring")
	}
}

// Test the duration() method
func Test_duration(t *testing.T) {

	// Create two times with a known difference (one hour)
	ping1 := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	ping2 := ping1.Add(time.Hour)

	// Create a peerAddr
	p := newPeerAddr(nil, ping1)

	// Ping the address five more times
	for i := 0; i < 5; i++ {
		p.ping(ping2)
	}

	// The duration should be one hour
	if p.duration() != time.Hour {
		t.Fatal("Duration does not match")
	}
}

// Test the isExpired() method
func Test_isExpired(t *testing.T) {

	// Create two times with a known difference (five seconds)
	time1 := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	time2 := time1.Add(5 * time.Second)

	// Create a new peerAddr with the first time
	p := newPeerAddr(nil, time1)

	// Assuming a timeout of three seconds, the address should have expired
	if !p.isExpired(3*time.Second, time2) {
		t.Fatal("Address should have expired")
	}
}
