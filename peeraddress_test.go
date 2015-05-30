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

// Test the Ping() method
func Test_Ping(t *testing.T) {

	// Create a new peerAddress and confirm that it contains one item
	p := newPeerAddress(nil, time.Now())
	if validElementsInRing(p.lastPing) != 1 {
		t.Fatal("Expected one element in ring")
	}

	// Register a ping and confirm that the ring now contains two items
	p.Ping(time.Now())
	if validElementsInRing(p.lastPing) != 2 {
		t.Fatal("Expected two elements in ring")
	}
}

// Test the Duration() method
func Test_Duration(t *testing.T) {

	// Create two times with a known difference (one hour)
	ping1 := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	ping2 := ping1.Add(time.Hour)

	// Create a peerAddress
	p := newPeerAddress(nil, ping1)

	// Ping the address five more times
	for i := 0; i < 5; i++ {
		p.Ping(ping2)
	}

	// The duration should be one hour
	if p.Duration() != time.Hour {
		t.Fatal("Duration does not match")
	}
}
