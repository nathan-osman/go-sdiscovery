package sdiscovery

import (
	"log"
	"net"
	"sync"
	"time"
)

// Network interface manager and packet monitor
type communicator struct {
	packetReceived chan *packet
	sendChan       chan []byte
	stopChan       chan struct{}
	connections    map[string][]*connection
	port           int
}

// Create a new communicator
func newCommunicator(pollInterval time.Duration, port int) *communicator {

	// Create the communicator, including the channel that will be used
	// for receiving the individual packets and the mapping between
	// interface names and the individual connections for them
	c := &communicator{
		packetReceived: make(chan *packet),
		sendChan:       make(chan []byte),
		stopChan:       make(chan struct{}),
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
	defer monitor.stop()

	// Create a WaitGroup for each of the sockets so that we can
	// ensure all of them end before closing the packet channel
	var waitGroup sync.WaitGroup

loop:
	for {
		select {
		case name := <-monitor.interfaceAdded:
			c.addInterface(name, &waitGroup)
		case name := <-monitor.interfaceRemoved:
			c.removeInterface(name)
		case data := <-c.sendChan:

			// Send on each of the connections
			for _, connections := range c.connections {
				for _, connection := range connections {
					connection.send(data)
				}
			}

		case <-c.stopChan:
			break loop
		}
	}

	// Stop all of the connections
	for name, _ := range c.connections {
		c.removeInterface(name)
	}

	// Wait for the connections to finish then close the channel
	waitGroup.Wait()
	close(c.packetReceived)
}

// Add connections for the specified interface
func (c *communicator) addInterface(name string, waitGroup *sync.WaitGroup) {

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
			if conn, err := newConnection(c.packetReceived, waitGroup, ifi, c.port, multicast); err != nil {
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
			connection.stop()
		}

		// Remove the item from the map
		delete(c.connections, name)
	}
}

// Send the specified data on each of the connections
func (c *communicator) send(data []byte) {
	c.sendChan <- data
}

// Stop the goroutine by closing the channels
func (c *communicator) stop() {
	close(c.sendChan)
	close(c.stopChan)
}
