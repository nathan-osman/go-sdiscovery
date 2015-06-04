package util

import (
	"bytes"
	"net"
	"testing"
)

// Test the ability to convert a CIDR to a broadcast address.
func Test_BroadcastIPFromCIDR(t *testing.T) {

	var (
		ip  net.IP
		err error
	)

	// Obtain the broadcast address for the provided CIDR.
	if ip, err = BroadcastIPFromCIDR("192.168.1.1/24"); err != nil {
		t.Fatal(err)
	}

	// Ensure that the IP addresses match.
	if !bytes.Equal(ip, net.IP{192, 168, 1, 255}) {
		t.Fatal("IP addresses do not match")
	}

	// Ensure that an error is generated for an IPv6 address.
	if ip, err = BroadcastIPFromCIDR("::1/128"); err == nil {
		t.Fatal("Expected error for IPv6 address")
	}
}
