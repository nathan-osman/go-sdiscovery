package util

import (
	"bytes"
	"net"
	"testing"
)

// Ensure that the broadcast address is correctly derived from a CIDR.
func Test_BroadcastIPFromCIDR(t *testing.T) {

	// Obtain the broadcast address for the provided CIDR.
	if ip, err := BroadcastIPFromCIDR("192.168.1.1/24"); err != nil {
		t.Fatal(err)
	} else {

		// Ensure that the IP addresses match.
		if !bytes.Equal(ip, net.IP{192, 168, 1, 255}) {
			t.Fatal("IP addresses do not match")
		}
	}

	// Ensure that an error is generated for an IPv6 address.
	if _, err := BroadcastIPFromCIDR("::1/128"); err == nil {
		t.Fatal("Expected error for IPv6 address")
	}
}
