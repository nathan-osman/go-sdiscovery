package sdiscovery

import (
	"log"
	"net"
	"time"
)

// Network interface manager and packet monitor
type communicator struct {
	PacketReceived chan packet
	stop           chan struct{}
	monitor        *monitor
	connections    map[string][]*connection
	port           int
}

// Create a new communicator
func newCommunicator(pollInterval time.Duration, port int) *communicator {

	// Create the communicator, including the channel that will be used
	// for receiving the individual packets and stopping the communicator
	c := &communicator{
		PacketReceived: make(chan packet),
		stop:           make(chan struct{}),
		monitor:        newMonitor(pollInterval),
		connections:    make(map[string][]*connection),
		port:           port,
	}

	// Spawn a goroutine that deals with connections
	go c.run()

	return c
}

// Add and remove connections as interfaces are added and removed
func (c *communicator) run() {
	for {
		select {
		case name := <-c.monitor.InterfaceAdded:
			c.addInterface(name)
		case name := <-c.monitor.InterfaceRemoved:
			c.removeInterface(name)
		case <-c.stop:
			return
		}
	}
}

// Add connections for the specified interface
func (c *communicator) addInterface(name string) {

	connections := make([]*connection, 2)

	// Attempt to find the interface by name
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		log.Println("[ERR]", err)
		return
	}

	// Add a connection for broadcast and multicast if supported
	for _, multicast := range []bool{true, false} {
		if ifi.Flags&net.FlagBroadcast != 0 {
			if conn, err := newConnection(c.PacketReceived, ifi, c.port, multicast); err != nil {
				log.Println("[ERR]", err)
			} else {
				connections = append(connections, conn)
			}
		}
	}

	// Create a new entry in the map for the connections
	if len(connections) != 0 {
		c.connections[name] = connections
	}
}

// Remove all connections for the specified interface
func (c *communicator) removeInterface(name string) {

	// Check if the interface exists
	if connections, exists := c.connections[name]; exists {

		// Stop the connections
		for _, connection := range connections {
			connection.Stop()
		}

		// Remove the item from the map
		delete(c.connections, name)
	}
}

// Send data on each of the connections
func (c *communicator) Send(data []byte) {
	for _, connections := range c.connections {
		for _, connection := range connections {
			connection.Send(data)
		}
	}
}

// Stop and close all of the connections
func (c *communicator) Stop() {

	// Stop the goroutine monitoring the interfaces
	close(c.stop)

	// Stop the monitor
	c.monitor.Stop()

	// Stop each of the individual connections
	for name, _ := range c.connections {
		c.removeInterface(name)
	}
}
