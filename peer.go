package sdiscovery

import (
	"encoding/json"
	"net"
	"time"
)

// Peer address obtained from received packets
type PeerAddress struct {
	IP         *net.IP
	lastPacket time.Time
}

// Peer found through local discovery
type Peer struct {
	ID        string          `json:"id"`
	UserData  json.RawMessage `json:"user_data"`
	Addresses []PeerAddress
}

// Check each of the peer's addresses to determine if any have expired. The
// specified timeout value is used to determine when this occurs.
func (p *Peer) Update(timeout time.Duration) {

	// Obtain the current time
	curTime := time.Now()

	// Create an empty slice pointing to the old array and filter the
	// addresses based on whether they have expired or not
	addresses := p.Addresses[:0]
	for _, addr := range p.Addresses {
		if !addr.lastPacket.Add(timeout).Before(curTime) {
			addresses = append(addresses, addr)
		}
	}
	p.Addresses = addresses
}

// Determine if the peer has any valid addresses remaining.
func (p *Peer) HasExpired() bool {
	return len(p.Addresses) == 0
}
