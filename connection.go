package sdiscovery

import (
	"encoding/binary"
	"errors"
	"net"
)

// Send and receive packets over a network connection
type connection struct {
	conn      *net.UDPConn
	bcastAddr *net.Addr
}

// Derive the broadcast address from an address in CIDR notation
func broadcastAddress(s string) (net.IP, error) {

	// Parse the CIDR
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}

	// Ensure that this is an IPv4 address
	if ipnet.IP.To4() == nil {
		return nil, errors.New("Not an IPv4 address")
	}

	// Convert the IP address and mask to 32-bit integers
	// Note that byte order is actually irrelevant here
	ip := binary.LittleEndian.Uint32(ipnet.IP)
	mask := binary.LittleEndian.Uint32(ipnet.Mask)

	// Calculate the broadcast address
	addr := make([]byte, 4)
	binary.LittleEndian.PutUint32(addr, ip&mask|^mask)

	return addr, nil
}

// Create a new connection
func newConnection(ifi *net.Interface, port int, multicast bool) (*connection, error) {

	var (
		conn      *net.UDPConn
		bcastAddr *net.Addr
		err       error
	)

	if multicast {

		// Connect to the all nodes link-local IPv6 address
		conn, err = net.ListenMulticastUDP("udp6", ifi, &net.UDPAddr{
			IP:   net.IPv6linklocalallnodes,
			Port: port,
		})

	} else {

		// Obtain all of the interface addresses
		addrs, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}

		// Attempt to find an IPv4 broadcast address
		var ip net.IP
		for _, addr := range addrs {
			if ip, err = broadcastAddress(addr.String()); err != nil {
				continue
			}
		}

		// Error if no addresses were found
		if ip == nil {
			return nil, errors.New("Unable to find a broadcast address")
		}

		bcastAddr := &net.UDPAddr{
			IP:   ip,
			Port: port,
		}
		conn, err = net.ListenUDP("udp4", bcastAddr)
	}

	// Check for an error
	if err != nil {
		return nil, err
	}

	return &connection{
		conn:      conn,
		bcastAddr: bcastAddr,
	}, nil
}
