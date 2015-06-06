package comm

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/nathan-osman/go-sdiscovery/util"
)

// Manages connections on available network interfaces.
type Communicator struct {
	PacketChan  chan *Packet
	sendChan    chan *Packet
	connections map[string][]*connection
	port        int
}

// Return a list of interface names.
func interfaceNames() (util.StrMap, error) {

	// Retrieve the current list of interfaces.
	ifis, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Create a map of the interface names.
	newNames := make(util.StrMap)
	for _, ifi := range ifis {
		newNames[ifi.Name] = nil
	}

	return newNames, nil
}

// Create a new communicator.
func NewCommunicator(pollInterval time.Duration, port int) *Communicator {

	// Create the communicator, including the channel that will be used
	// for receiving the individual packets.
	c := &Communicator{
		PacketChan:  make(chan *Packet),
		sendChan:    make(chan *Packet),
		connections: make(map[string][]*connection),
		port:        port,
	}

	// Spawn a goroutine that manages connections.
	go c.run(pollInterval)

	return c
}

// Add and remove connections as interfaces are added and removed.
func (c *Communicator) run(pollInterval time.Duration) {

	// Enumerate interface additions and removals.
	ticker := time.NewTicker(pollInterval)
	ifiEnum := util.NewStrEnum(ticker.C, interfaceNames)
	defer ticker.Stop()

	// Create a WaitGroup for each of the sockets so that we can ensure all of
	// them end before closing the packet channel.
	var waitGroup sync.WaitGroup

loop:
	for {
		select {
		case name := <-ifiEnum.StringAdded:
			c.addInterface(name, &waitGroup)
		case name := <-ifiEnum.StringRemoved:
			c.removeInterface(name)
		case data, ok := <-c.sendChan:

			// If the receive was successful, send the packet on each of the
			// connections. Otherwise, quit the loop.
			if ok {
				for _, connections := range c.connections {
					for _, connection := range connections {
						connection.send(data)
					}
				}
			} else {
				break loop
			}
		}
	}

	// Stop all of the connections.
	for name, _ := range c.connections {
		c.removeInterface(name)
	}

	// Wait for the connections to finish then close the channel.
	waitGroup.Wait()
	close(c.PacketChan)
}

// Add connections for the specified interface.
func (c *Communicator) addInterface(name string, waitGroup *sync.WaitGroup) {

	// Assume that most interfaces will have at most two addresses.
	connections := make([]*connection, 0, 2)

	// Attempt to find the interface by name.
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		log.Println("[ERR]", err)
		return
	}

	// Add a connection for broadcast and multicast addresses if present.
	if ifi.Flags&net.FlagMulticast != 0 {
		if conn, err := newConnection(c.PacketChan, waitGroup, ifi, c.port, multicast); err != nil {
			log.Println("[ERR]", err)
		} else {
			connections = append(connections, conn)
		}
	}
	if ifi.Flags&net.FlagBroadcast != 0 {
		if conn, err := newConnection(c.PacketChan, waitGroup, ifi, c.port, broadcast); err != nil {
			log.Println("[ERR]", err)
		} else {
			connections = append(connections, conn)
		}
	}

	// Create a new entry in the map for the connections (if any).
	if len(connections) != 0 {
		c.connections[name] = connections
	}
}

// Remove all connections for the specified interface.
func (c *Communicator) removeInterface(name string) {

	// Check if the interface exists.
	if connections, ok := c.connections[name]; ok {

		// Stop the connections.
		for _, connection := range connections {
			connection.stop()
		}

		// Remove the item from the map.
		delete(c.connections, name)
	}
}

// Send the specified packet on each of the connections.
func (c *Communicator) Send(pkt *Packet) {
	c.sendChan <- pkt
}

// Stop the goroutine by closing the send channel.
func (c *Communicator) Stop() {
	close(c.sendChan)
}
