package sdiscovery

import (
	"log"
	"net"
	"time"
)

// Monitor attached network interfaces
type monitor struct {
	InterfaceAdded   chan string
	InterfaceRemoved chan string
	stop             chan struct{}
}

// Create a new interface monitor
func newMonitor(pollInterval time.Duration) *monitor {

	// Create the monitor
	m := &monitor{
		InterfaceAdded:   make(chan string),
		InterfaceRemoved: make(chan string),
		stop:             make(chan struct{}),
	}

	// Spawn a goroutine to perform the enumeration at the scheduled interval
	go m.run(pollInterval)

	return m
}

// Regularly poll for new network interfaces
func (m *monitor) run(pollInterval time.Duration) {

	// Create a map to store the interface names between enumerations
	var oldNames map[string]struct{} = m.enumerate(map[string]struct{}{})

	for {
		select {
		case <-time.After(pollInterval):
			oldNames = m.enumerate(oldNames)
		case <-m.stop:
			close(m.InterfaceAdded)
			close(m.InterfaceRemoved)
			return
		}
	}
}

// Check for changes to the list of network interfaces
func (m *monitor) enumerate(oldNames map[string]struct{}) map[string]struct{} {

	// Retrieve the current list of interfaces
	ifis, err := net.Interfaces()
	if err != nil {

		// Assume that this error is temporary and try again next time
		log.Println("[ERR]", err)
		return oldNames
	}

	// Create a map of the interface names
	newNames := make(map[string]struct{})
	for _, ifi := range ifis {
		newNames[ifi.Name] = struct{}{}
	}

	// Compare the two maps
	m.compare(oldNames, newNames, m.InterfaceAdded)
	m.compare(newNames, oldNames, m.InterfaceRemoved)

	return newNames
}

// Compare two maps and notify of any changes on the specified channel
func (m *monitor) compare(a, b map[string]struct{}, notify chan<- string) {
	for name, _ := range b {
		if _, exists := a[name]; !exists {
			select {
			case notify <- name:
			case <-m.stop:
			}
		}
	}
}

// Stop monitoring for new network interfaces
func (m *monitor) Stop() {
	close(m.stop)
}
