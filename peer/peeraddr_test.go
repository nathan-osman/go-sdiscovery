package peer

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

// Generate times that differ by the specified durations
func generateTimes(durations ...time.Duration) []time.Time {

	// Create a slice for storing the times and add the first one
	times := make([]time.Time, len(durations)+1)
	times[0] = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Add each of the durations
	for i, duration := range durations {
		times[i+1] = times[i].Add(duration)
	}

	return times
}

// Test the ping() method
func Test_peerAddr_ping(t *testing.T) {

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
func Test_peerAddr_duration(t *testing.T) {

	// Create two times with a known difference (one second)
	times := generateTimes(time.Second)

	// Create a peerAddr
	p := newPeerAddr(nil, times[0])

	// Ping the address five more times
	for i := 0; i < 5; i++ {
		p.ping(times[1])
	}

	// The duration should be one hour
	if p.duration() != time.Second {
		t.Fatal("Duration does not match")
	}
}

// Test the isExpired() method
func Test_peerAddr_isExpired(t *testing.T) {

	// Create two times with a known difference (two seconds)
	times := generateTimes(2 * time.Second)

	// Create a new peerAddr with the first time
	p := newPeerAddr(nil, times[0])

	// Assuming a timeout of one second, the address should have expired
	if !p.isExpired(1*time.Second, times[1]) {
		t.Fatal("Address should have expired")
	}
}
