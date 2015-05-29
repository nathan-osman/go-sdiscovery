package sdiscovery

import (
	"net"
	"sync"
)

// Data about an individual packet received from a connection
type packet struct {
	IP   net.IP
	Data []byte
}

// Sender and receiver for packets over a network connection
type connection struct {
	packetReceived chan<- packet
	stop           chan struct{}
	conn           *net.UDPConn
	waitGroup      *sync.WaitGroup
}

// Create a new connection for sending and receiving packets
func newConnection(packetReceived chan<- packet, waitGroup *sync.WaitGroup, multicast bool, ifi *net.Interface, port int) (*connection, error) {

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
		ip, err := findBroadcastAddress(ifi)
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
		packetReceived: packetReceived,
		stop:           make(chan struct{}),
		conn:           conn,
		waitGroup:      waitGroup,
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
		b := make([]byte, 1000)

		// Read the packet, quitting the goroutine on error
		n, addr, err := c.conn.ReadFromUDP(b)
		if err != nil {
			return
		}

		// Create the packet
		p := packet{
			IP:   addr.IP,
			Data: b[:n],
		}

		// Write the packet to the channel
		select {
		case c.packetReceived <- p:
		case <-c.stop:
		}
	}
}

// Send a packet over the connection
func (c *connection) Send(data []byte) error {
	_, err := c.conn.WriteToUDP(data, c.conn.LocalAddr().(*net.UDPAddr))
	return err
}

// Stop listening for incoming packets
func (c *connection) Stop() {
	c.conn.Close()
	close(c.stop)
}
