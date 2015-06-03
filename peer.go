package sdiscovery

import (
	"time"
)

// peer maintains information about a peer discovered on the network. Because
// the struct may be used from multiple goroutines, all access to members must
// be done through accessors that lock a mutex.
type peer struct {
	userData []byte
	addrs    []*peerAddr
}

// Record a ping from the specified address
func (p *peer) ping(pkt *packet, curTime time.Time) {

	// Store userData
	p.userData = pkt.userData

	// Attempt to find a matching address
	for _, addr := range p.addrs {
		if pkt.ip == addr.ip {
			addr.ping(curTime)
			return
		}
	}

	// No matching address was found, add a new one
	p.addrs = append(p.addrs, newPeerAddr(pkt.ip, curTime))
}

// Remove all expired addresses and sort those that remain
func (p *peer) isExpired(timeout time.Duration, curTime time.Time) bool {

	// Create an empty slice pointing to the old array and filter the
	// addresses based on whether they have expired or not
	addrs := p.addrs[:0]
	for _, addr := range p.addrs {
		if !addr.isExpired(timeout, curTime) {
			addrs = append(addrs, addr)
		}
	}
	p.addrs = addrs

	// Return true if no addresses remain
	return len(p.addrs) == 0
}
