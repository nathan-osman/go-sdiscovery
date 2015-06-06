package comm

import (
	"encoding/json"
	"net"
)

// Packet represents an individual packet received from a network interface.
type Packet struct {
	IP       net.IP `json:"-"`         // IP address from which the packet was obtained
	ID       string `json:"id"`        // ID of the peer that sent the packet
	UserData []byte `json:"user_data"` // custom data provided by the peer
}

// Create a new packet using the specified IP address and JSON data.
func NewPacketFromJSON(ip net.IP, data []byte) (*Packet, error) {

	// Create the packet with the provided IP address.
	pkt := &Packet{
		IP: ip,
	}

	// Attempt to unmarshal the JSON into the packet.
	if err := json.Unmarshal(data, pkt); err != nil {
		return nil, err
	}

	return pkt, nil
}

// Convert the packet to JSON.
func (p *Packet) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}
