package sdiscovery

import (
	"log"
	"net"
	"time"
)

// Network interface manager and packet monitor
type communicator struct {
	PacketReceived chan packet
	send           chan []byte
	stop           chan struct{}
	connections    map[string][]*connection
	port           int
}

// Create a new communicator
func newCommunicator(pollInterval time.Duration, port int) *communicator {

	// Create the communicator, including the channel that will be used
	// for receiving the individual packets and the mapping between
	// interface names and the individual connections for them
	c := &communicator{
		PacketReceived: make(chan packet),
		send:           make(chan []byte),
		stop:           make(chan struct{}),
		connections:    make(map[string][]*connection),
		port:           port,
	}

	// Spawn a goroutine that manages connections
	go c.run(pollInterval)

	return c
}

// Add and remove connections as interfaces are added and removed
func (c *communicator) run(pollInterval time.Duration) {

	// Monitor for interface additions and removals
	monitor := newMonitor(pollInterval)
	defer monitor.Stop()

loop:
	for {
		select {
		case name := <-monitor.InterfaceAdded:
			c.addInterface(name)
		case name := <-monitor.InterfaceRemoved:
			c.removeInterface(name)
		case data := <-c.send:

			// Send on each of the connections
			for _, connections := range c.connections {
				for _, connection := range connections {
					connection.Send(data)
				}
			}

		case <-c.stop:
			break loop
		}
	}

	// Stop all of the connections
	for name, _ := range c.connections {
		c.removeInterface(name)
	}

	// TODO: this could be a problem if a connection attempts
	// to write to the channel after we close it here
	close(c.PacketReceived)
}

// Add connections for the specified interface
func (c *communicator) addInterface(name string) {

	// Assume that most interfaces will have at most two addresses
	connections := make([]*connection, 0, 2)

	// Attempt to find the interface by name
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		log.Println("[ERR]", err)
		return
	}

	// Add a connection for broadcast and multicast addresses if present
	for _, multicast := range []bool{true, false} {
		if ifi.Flags&net.FlagBroadcast != 0 {
			if conn, err := newConnection(c.PacketReceived, ifi, c.port, multicast); err != nil {
				log.Println("[ERR]", err)
			} else {
				connections = append(connections, conn)
			}
		}
	}

	// Create a new entry in the map for the connections (if any)
	if len(connections) != 0 {
		c.connections[name] = connections
	}
}

// Remove all connections for the specified interface
func (c *communicator) removeInterface(name string) {

	// Check if the interface exists
	if connections, ok := c.connections[name]; ok {

		// Stop the connections
		for _, connection := range connections {
			connection.Stop()
		}

		// Remove the item from the map
		delete(c.connections, name)
	}
}

// Send the specified data on each of the connections
func (c *communicator) Send(data []byte) {
	c.send <- data
}

// Stop the goroutine by closing the channels
func (c *communicator) Stop() {
	close(c.send)
	close(c.stop)
}
