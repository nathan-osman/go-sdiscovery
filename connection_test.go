package sdiscovery

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// Test the broadcastAddress function
func Test_broadcastAddressFromCIDR(t *testing.T) {

	var ip net.IP
	var err error

	// Obtain the broadcast address for the provided CIDR
	if ip, err = broadcastAddressFromCIDR("192.168.1.1/24"); err != nil {
		t.Fatal(err)
	}

	// Ensure that the IP addresses match
	if !bytes.Equal(ip, net.IP{192, 168, 1, 255}) {
		t.Fatal("IP addresses do not match")
	}

	// Ensure that an error is generated for an IPv6 address
	if ip, err = broadcastAddressFromCIDR("::1/128"); err == nil {
		t.Fatal("Expected error for IPv6 address")
	}
}

// Testing the findBroadcastAddress function is virtually impossible since
// there is no way (AFAIK) to simulate an interface for testing

// Attempt to find an interface with the specified flag
func findInterfaceWithFlag(flag net.Flags) (*net.Interface, error) {

	// Obtain the list of interfaces
	ifis, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Return the first one that matches
	for _, ifi := range ifis {
		if ifi.Flags&flag != 0 {
			return &ifi, nil
		}
	}

	// None matched - return nil
	return nil, nil
}

// Test that packets are correctly sent and received via broadcast
func Test_connection_broadcast(t *testing.T) {

	// Attempt to find a broadcast interface
	ifi, err := findInterfaceWithFlag(net.FlagBroadcast)
	if err != nil {
		t.Fatal(err)
	}

	// Skip the test if none was found
	if ifi == nil {
		t.Skip("No broadcast interface found")
	}

	// Create the connection with a randomly chosen port
	conn, err := newConnection(ifi, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	// Send a packet
	packet := []byte(`test`)
	if err := conn.Send(packet); err != nil {
		t.Fatal(err)
	}

	// Receive the packet
	select {
	case b := <-conn.PacketReceived:
		if !bytes.Equal(b, packet) {
			t.Fatal("Packet contents do not match")
		}
	case <-time.NewTicker(50 * time.Millisecond).C:
		t.Fatal("Timeout waiting for broadcast packet")
	}
}

// Test that packets are correctly sent and received via multicast
func Test_connection_multicast(t *testing.T) {
	//...
}
