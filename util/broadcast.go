package util

import (
	"encoding/binary"
	"errors"
	"net"
)

// Derive the broadcast IP address from an IP address in CIDR notation.
func BroadcastIPFromCIDR(cidr string) (net.IP, error) {

	// Parse the CIDR.
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	// Ensure that it is a valid IPv4 address.
	if ipnet.IP.To4() == nil {
		return nil, errors.New("Not an IPv4 address")
	}

	// Convert the IP address and mask to 32-bit integers. Note that byte order
	// is actually irrelevant since calculations are done at the bit level.
	ip := binary.LittleEndian.Uint32(ipnet.IP)
	mask := binary.LittleEndian.Uint32(ipnet.Mask)

	// Calculate the broadcast address.
	bcastIP := make([]byte, 4)
	binary.LittleEndian.PutUint32(bcastIP, ip&mask|^mask)

	return bcastIP, nil
}

// Find a broadcast address for the provided network interface. If the
// interface does not contain an IPv4 address, this function will fail.
func FindBroadcastAddress(ifi *net.Interface) (net.IP, error) {

	// Obtain all of the addresses on the interface.
	addrs, err := ifi.Addrs()
	if err != nil {
		return nil, err
	}

	// For each of the addresses, check if a valid broadcast address exists.
	for _, addr := range addrs {
		if ip, err := BroadcastIPFromCIDR(addr.String()); err == nil {
			return ip, nil
		}
	}

	return nil, errors.New("No broadcast address was found")
}
