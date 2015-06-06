package peer

import (
	"net"
	"sort"
	"time"

	"github.com/nathan-osman/go-sdiscovery/comm"
)

type peerSlice []*peerAddr

// Peer maintains information about a peer discovered on the network. Because
// the struct may be used from multiple goroutines, all access to members must
// be done through accessors that lock a mutex.
type Peer struct {
	UserData []byte
	addrs    peerSlice
}

func (a peerSlice) Len() int           { return len(a) }
func (a peerSlice) Less(i, j int) bool { return a[i].duration() < a[j].duration() }
func (a peerSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Record a ping from the specified address.
func (p *Peer) Ping(pkt *comm.Packet, curTime time.Time) {

	// Store userData.
	p.UserData = pkt.UserData

	// Attempt to find a matching address.
	for _, addr := range p.addrs {
		if pkt.IP.Equal(addr.ip) {
			addr.ping(curTime)
			return
		}
	}

	// No matching address was found, add a new one.
	p.addrs = append(p.addrs, newPeerAddr(pkt.IP, curTime))
}

// Obtain a sorted list of all addresses for the peer.
func (p *Peer) Addrs() []net.IP {

	// First sort the addresses
	sort.Sort(p.addrs)

	// Build a slice of IP addresses
	ips := make([]net.IP, len(p.addrs))
	for i, addr := range p.addrs {
		ips[i] = addr.ip
	}

	return ips
}

// Remove all expired addresses and sort those that remain.
func (p *Peer) IsExpired(timeout time.Duration, curTime time.Time) bool {

	// Create an empty slice pointing to the old array and filter the
	// addresses based on whether they have expired or not.
	addrs := p.addrs[:0]
	for _, addr := range p.addrs {
		if !addr.isExpired(timeout, curTime) {
			addrs = append(addrs, addr)
		}
	}
	p.addrs = addrs

	// Return true if no addresses remain.
	return len(p.addrs) == 0
}
