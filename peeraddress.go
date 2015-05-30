package sdiscovery

import (
	"container/ring"
	"net"
	"time"
)

// peerAddress contains a single address that has received packets and a ring
// that keeps track of the time between the last few pings received. In this
// case, the lower the duration, the better.
type peerAddress struct {
	IP       *net.IP
	lastPing *ring.Ring
}

// Create a new peerAddress
func newPeerAddress(ip *net.IP, t time.Time) *peerAddress {

	// Create the new peer address
	p := &peerAddress{
		IP:       ip,
		lastPing: ring.New(6),
	}

	// Record the current ping
	p.lastPing.Value = t

	return p
}

// Register a ping against the address
func (p *peerAddress) Ping(t time.Time) {

	// Advance forward and record the current time
	p.lastPing = p.lastPing.Next()
	p.lastPing.Value = t
}

// Determine the duration between the oldest and most recent packet
func (p *peerAddress) Duration() time.Duration {

	// Attempt to convert the "next" value (the oldest) to time.Time
	oldestPing, _ := p.lastPing.Next().Value.(time.Time)

	// Return the difference between the two pings
	return p.lastPing.Value.(time.Time).Sub(oldestPing)
}
