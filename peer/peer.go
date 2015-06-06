package peer

import (
	"time"

	"github.com/nathan-osman/go-sdiscovery/conn"
)

// Peer maintains information about a peer discovered on the network. Because
// the struct may be used from multiple goroutines, all access to members must
// be done through accessors that lock a mutex.
type Peer struct {
	UserData []byte
	Addrs    []*peerAddr
}

// Record a ping from the specified address.
func (p *Peer) Ping(pkt *conn.Packet, curTime time.Time) {

	// Store userData.
	p.UserData = pkt.UserData

	// Attempt to find a matching address.
	for _, addr := range p.Addrs {
		if pkt.IP.Equal(addr.ip) {
			addr.ping(curTime)
			return
		}
	}

	// No matching address was found, add a new one.
	p.Addrs = append(p.Addrs, newPeerAddr(pkt.IP, curTime))
}

// Remove all expired addresses and sort those that remain.
func (p *Peer) IsExpired(timeout time.Duration, curTime time.Time) bool {

	// Create an empty slice pointing to the old array and filter the
	// addresses based on whether they have expired or not.
	addrs := p.Addrs[:0]
	for _, addr := range p.Addrs {
		if !addr.isExpired(timeout, curTime) {
			addrs = append(addrs, addr)
		}
	}
	p.Addrs = addrs

	// Return true if no addresses remain.
	return len(p.Addrs) == 0
}
