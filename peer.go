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

// Check the peer's addresses for expiry and determine if the peer itself
// has expired based on the last ping received.
func (p *Peer) HasExpired(timeout time.Duration) bool {

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

	// If all of the addresses have expired, then so has the device
	return len(p.Addresses) == 0
}
