package comm

import (
	"net"
	"sync"

	"github.com/nathan-osman/go-sdiscovery/util"
)

type packetType int

const (
	multicast packetType = iota
	broadcast
)

// connection provides methods for sending and receiving packets from a
// specific address on a network interface.
type connection struct {
	stopChan chan interface{}
	conn     *net.UDPConn
}

// Create a new multicast (IPv6) connection to the specified interface.
func multicastConnection(ifi *net.Interface, port int) (*net.UDPConn, error) {

	// Use the all nodes link-local IPv6 address.
	return net.ListenMulticastUDP("udp6", ifi, &net.UDPAddr{
		IP:   net.IPv6linklocalallnodes,
		Port: port,
	})
}

// Create a new broadcast (IPv4) connection to the specified interface.
func broadcastConnection(ifi *net.Interface, port int) (*net.UDPConn, error) {

	// Attempt to find an IPv4 broadcast address.
	ip, err := util.FindBroadcastAddress(ifi)
	if err != nil {
		return nil, err
	}

	// Use the broadcast address that was found.
	return net.ListenUDP("udp4", &net.UDPAddr{
		IP:   ip,
		Port: port,
	})
}

// Create a new connection for sending and receiving packets.
func newConnection(packetChan chan<- *Packet, waitGroup *sync.WaitGroup, ifi *net.Interface, port int, pType packetType) (*connection, error) {

	var (
		conn *net.UDPConn
		err  error
	)

	// Use the appropriate initializer.
	switch pType {
	case multicast:
		conn, err = multicastConnection(ifi, port)
	case broadcast:
		conn, err = broadcastConnection(ifi, port)
	}

	// Check for an error.
	if err != nil {
		return nil, err
	}

	// Create the connection.
	c := &connection{
		stopChan: make(chan interface{}),
		conn:     conn,
	}

	// Spawn a goroutine to read from the socket.
	go c.run(packetChan, waitGroup)

	return c, nil
}

// Continuously read packets from the connection.
func (c *connection) run(packetChan chan<- *Packet, waitGroup *sync.WaitGroup) {

	// Ensure that the WaitGroup is properly updated.
	waitGroup.Add(1)
	defer waitGroup.Done()

loop:
	for {

		// Put a hard cap of 1000 bytes on the packet size.
		b := make([]byte, 1000)

		// Read the packet, quitting on error.
		n, addr, err := c.conn.ReadFromUDP(b)
		if err != nil {
			break
		}

		// Attempt to create the packet.
		pkt, err := NewPacketFromJSON(addr.IP, b[:n])
		if err != nil {
			continue
		}

		// Write the packet to the channel.
		select {
		case packetChan <- pkt:
		case <-c.stopChan:
			break loop
		}
	}
}

// Send a packet.
func (c *connection) send(pkt *Packet) error {

	// Convert the packet to JSON.
	data, err := pkt.ToJSON()
	if err != nil {
		return err
	}

	// Write the packet to the connection.
	_, err = c.conn.WriteToUDP(data, c.conn.LocalAddr().(*net.UDPAddr))
	return err
}

// Stop listening for incoming packets.
func (c *connection) stop() {
	c.conn.Close()
	close(c.stopChan)
}
