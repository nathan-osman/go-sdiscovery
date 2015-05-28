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
	oldNames         map[string]struct{}
}

// Create a new monitor
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

	// Immediately enumerate the interfaces
	m.enumerate()

	// Create a ticker to schedule interface enumeration
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.enumerate()
		case <-m.stop:
			close(m.InterfaceAdded)
			close(m.InterfaceRemoved)
			return
		}
	}
}

// Check for changes to the list of network interfaces
func (m *monitor) enumerate() {

	// Retrieve the current list of interfaces
	ifis, err := net.Interfaces()
	if err != nil {

		// Assume that this error is temporary and try again next time
		log.Println("[ERR]", err)
		return
	}

	// Create a map of the interface names
	newNames := make(map[string]struct{})
	for _, ifi := range ifis {
		newNames[ifi.Name] = struct{}{}
	}

	// Write each of the new names to InterfaceAdded
	for name, _ := range newNames {
		if _, exists := m.oldNames[name]; !exists {
			select {
			case m.InterfaceAdded <- name:
			case <-m.stop:
			}
		}
	}

	// Write each of the missing names to InterfaceRemoved
	for name, _ := range m.oldNames {
		if _, exists := newNames[name]; !exists {
			select {
			case m.InterfaceRemoved <- name:
			case <-m.stop:
			}
		}
	}

	// Assign the new list to the monitor
	m.oldNames = newNames
}

// Stop monitoring for new network interfaces
func (m *monitor) Stop() {
	close(m.stop)
}
