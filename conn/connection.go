package conn

import (
	"net"
	"sync"

	"github.com/nathan-osman/go-sdiscovery/util"
)

// Sender and receiver for packets over a network connection
type connection struct {
	packetChan chan<- *packet
	stopChan   chan struct{}
	conn       *net.UDPConn
	waitGroup  *sync.WaitGroup
}

// Create a new connection for sending and receiving packets
func newConnection(packetChan chan<- *packet, waitGroup *sync.WaitGroup, ifi *net.Interface, port int, multicast bool) (*connection, error) {

	var (
		conn *net.UDPConn
		err  error
	)

	if multicast {

		// Use the all nodes link-local IPv6 address
		conn, err = net.ListenMulticastUDP("udp6", ifi, &net.UDPAddr{
			IP:   net.IPv6linklocalallnodes,
			Port: port,
		})

	} else {

		// Attempt to find an IPv4 broadcast address
		ip, err := util.FindBroadcastAddress(ifi)
		if err != nil {
			return nil, err
		}

		// Build the broadcast address
		conn, err = net.ListenUDP("udp4", &net.UDPAddr{
			IP:   ip,
			Port: port,
		})
	}

	// Check for an error
	if err != nil {
		return nil, err
	}

	// Create the connection
	c := &connection{
		packetChan: packetChan,
		stopChan:   make(chan struct{}),
		conn:       conn,
		waitGroup:  waitGroup,
	}

	// Spawn a goroutine to read from the socket
	go c.run()

	return c, nil
}

// Continuously read packets from the connection
func (c *connection) run() {

	// Ensure that the WaitGroup is properly updated
	c.waitGroup.Add(1)
	defer c.waitGroup.Done()

	for {

		// Put a hard cap of 1000 bytes on the packet size
		// TODO: grab the MTU from the interface
		b := make([]byte, 1000)

		// Read the packet, quitting on error
		n, addr, err := c.conn.ReadFromUDP(b)
		if err != nil {
			break
		}

		// Attempt to create the packet
		pkt, err := newPacketFromJSON(addr.IP, b[:n])
		if err != nil {
			continue
		}

		// Write the packet to the channel
		select {
		case c.packetChan <- pkt:
		case <-c.stopChan:
		}
	}
}

// Send a packet over the connection
func (c *connection) send(data []byte) error {
	_, err := c.conn.WriteToUDP(data, c.conn.LocalAddr().(*net.UDPAddr))
	return err
}

// Stop listening for incoming packets
func (c *connection) stop() {
	c.conn.Close()
	close(c.stopChan)
}
