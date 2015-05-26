package sdiscovery

import (
	"net"
	"reflect"
	"sort"
	"testing"
	"time"
)

// Test for the correct interface names being sent over the InterfaceAdded socket
func TestMonitor_InterfaceAdded(t *testing.T) {

	var foundNames []string
	var monitorNames []string

	// Obtain a list of all network interfaces
	ifis, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	// Put the names into the list
	for _, ifi := range ifis {
		foundNames = append(foundNames, ifi.Name)
	}

	// Create a monitor
	monitor := newMonitor(50 * time.Millisecond)

	// Wait for interfaces to be enumerated
	timeout := time.After(150 * time.Millisecond)

loop:
	for {
		select {
		case name := <-monitor.InterfaceAdded:
			monitorNames = append(monitorNames, name)
		case <-timeout:
			break loop
		}
	}

	// Sort the lists for proper comparison
	sort.Strings(foundNames)
	sort.Strings(monitorNames)

	// Compare the list of interfaces
	if !reflect.DeepEqual(foundNames, monitorNames) {
		t.Fatal("Interface names do not match")
	}
}
