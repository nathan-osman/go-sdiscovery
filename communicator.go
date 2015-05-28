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
func newCommunicator(ifiInterval time.Duration, port int) *communicator {

	// Create the communicator, including the channel that will be used
	// for receiving the individual packets and stopping the communicator
	c := &communicator{
		PacketReceived: make(chan packet),
		stop:           make(chan struct{}),
		monitor:        newMonitor(ifiInterval),
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
			delete(c.connections, name)
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

// Stop and close all of the connections
func (c *communicator) Stop() {

	// Stop the goroutine monitoring the interfaces
	c.stop <- struct{}{}

	// Stop the monitor
	c.monitor.Stop()

	// Stop each of the individual connections
	for _, connections := range c.connections {
		for _, connection := range connections {
			connection.Stop()
		}
	}
}
