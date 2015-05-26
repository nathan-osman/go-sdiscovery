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
	oldNames         map[string]struct{}
}

// Create a new monitor
func newMonitor(duration time.Duration) *monitor {

	// Create the monitor
	m := &monitor{
		InterfaceAdded:   make(chan string),
		InterfaceRemoved: make(chan string),
	}

	// Create a ticker to schedule interface enumeration
	ticker := time.NewTicker(duration)

	// Spawn a new goroutine to enumerate the interfaces
	go func() {
		m.enumerate()
		for {
			m.enumerate()
			<-ticker.C
		}
	}()

	return m
}

// Repeatedly poll for new network interfaces
func (m *monitor) enumerate() {

	// Fetch the current list of interfaces
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
			m.InterfaceAdded <- name
		}
	}

	// Write each of the missing names to InterfaceRemoved
	for name, _ := range m.oldNames {
		if _, exists := newNames[name]; !exists {
			m.InterfaceRemoved <- name
		}
	}

	// Assign the new list to the monitor
	m.oldNames = newNames
}
