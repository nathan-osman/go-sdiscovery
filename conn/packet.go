package conn

import (
	"encoding/json"
	"net"
)

// Individual packet received from an interface
type packet struct {
	ip       net.IP
	ID       string `json:"id"`
	UserData []byte `json:"user_data"`
}

// Create a new packet from JSON data
func newPacketFromJSON(ip net.IP, data []byte) (*packet, error) {

	// Create the packet with the provided IP address
	pkt := &packet{
		ip: ip,
	}

	// Attempt to unmarshal the JSON into the packet
	if err := json.Unmarshal(data, &pkt); err != nil {
		return nil, err
	}

	return pkt, nil
}

// Convert a packet to JSON
func (p *packet) toJSON() ([]byte, error) {
	return json.Marshal(p)
}
