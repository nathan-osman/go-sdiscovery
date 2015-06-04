package sdiscovery

import (
	"log"
	"net"
	"time"
)

// Map of interface names
type interfaceMap map[string]interface{}

// Monitor attached network interfaces
type monitor struct {
	interfaceAdded   chan string
	interfaceRemoved chan string
	stopChan         chan interface{}
}

// Create a new interface monitor
func newMonitor(pollInterval time.Duration) *monitor {

	// Create the monitor
	m := &monitor{
		interfaceAdded:   make(chan string),
		interfaceRemoved: make(chan string),
		stopChan:         make(chan interface{}),
	}

	// Spawn a goroutine to perform the enumeration at the scheduled interval
	go m.run(pollInterval)

	return m
}

// Regularly poll for new network interfaces
func (m *monitor) run(pollInterval time.Duration) {

	// Create a map to store the interface names between enumerations
	var oldNames interfaceMap = m.enumerate(interfaceMap{})

	for {
		select {
		case <-time.After(pollInterval):
			oldNames = m.enumerate(oldNames)
		case <-m.stopChan:
			close(m.interfaceAdded)
			close(m.interfaceRemoved)
			return
		}
	}
}

// Check for changes to the list of network interfaces
func (m *monitor) enumerate(oldNames interfaceMap) interfaceMap {

	// Retrieve the current list of interfaces
	ifis, err := net.Interfaces()
	if err != nil {

		// Assume that this error is temporary and try again next time
		log.Println("[ERR]", err)
		return oldNames
	}

	// Create a map of the interface names
	newNames := make(interfaceMap)
	for _, ifi := range ifis {
		newNames[ifi.Name] = nil
	}

	// Compare the two maps
	m.compare(oldNames, newNames, m.interfaceAdded)
	m.compare(newNames, oldNames, m.interfaceRemoved)

	return newNames
}

// Compare two maps and notify of any changes on the specified channel
func (m *monitor) compare(a, b interfaceMap, notifyChan chan<- string) {
	for name, _ := range b {
		if _, exists := a[name]; !exists {
			select {
			case notifyChan <- name:
			case <-m.stopChan:
			}
		}
	}
}

// Stop monitoring for new network interfaces
func (m *monitor) stop() {
	close(m.stopChan)
}
