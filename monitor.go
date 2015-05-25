package sdiscovery

import (
	"log"
	"net"
	"time"
)

// Monitor attached network interfaces
type monitor struct {
	InterfaceAdded chan string
	oldNames       map[string]struct{}
	ticker         *time.Ticker
}

// Create a new monitor
func newMonitor(duration time.Duration) *monitor {

	// Create the monitor
	m := &monitor{
		InterfaceAdded: make(chan string),
		ticker:         time.NewTicker(duration),
	}

	// Spawn a new goroutine to monitor the interfaces
	go m.run()

	return m
}

// Repeatedly poll for new network interfaces
func (m *monitor) run() {

	for {

		// Fetch the current list of interfaces
		ifis, err := net.Interfaces()
		if err != nil {
			log.Println("[ERR]", err)
		}

		// Create a map of the interface names
		newNames := make(map[string]struct{})
		for _, ifi := range ifis {
			newNames[ifi.Name] = struct{}{}
		}

		// Write each of the new names to the channel
		for name, _ := range newNames {
			if _, exists := m.oldNames[name]; !exists {
				m.InterfaceAdded <- name
			}
		}

		// Wait for the next interval
		<-m.ticker.C
	}
}
