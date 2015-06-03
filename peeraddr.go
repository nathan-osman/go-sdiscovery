package sdiscovery

import (
	"container/ring"
	"net"
	"time"
)

// peerAddr contains a single address that has received packets and a ring
// that keeps track of the time between the last few pings received. In this
// case, the lower the duration, the better.
type peerAddr struct {
	ip       net.IP
	lastPing *ring.Ring
}

// Create a new peerAddr
func newPeerAddr(ip net.IP, curTime time.Time) *peerAddr {

	// Create the new peer address
	p := &peerAddr{
		ip:       ip,
		lastPing: ring.New(6),
	}

	// Record the current ping
	p.lastPing.Value = curTime

	return p
}

// Register a ping against the address
func (p *peerAddr) ping(curTime time.Time) {

	// Advance forward and record the current time
	p.lastPing = p.lastPing.Next()
	p.lastPing.Value = curTime
}

// Determine the duration between the oldest and most recent packet
func (p *peerAddr) duration() time.Duration {

	// Attempt to convert the "next" value (the oldest) to time.Time
	oldestPing, _ := p.lastPing.Next().Value.(time.Time)

	// Return the difference between the two pings
	return p.lastPing.Value.(time.Time).Sub(oldestPing)
}

// Determine if the address has exceeded the specified timeout
func (p *peerAddr) isExpired(timeout time.Duration, curTime time.Time) bool {
	return curTime.Sub(p.lastPing.Value.(time.Time)) >= timeout
}
