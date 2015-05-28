package sdiscovery

import (
	"net"
	"reflect"
	"sort"
	"testing"
	"time"
)

// Test that interface names are correctly enumerated
func Test_enumerate(t *testing.T) {

	// Obtain a list of all network interfaces
	ifis, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	// Create a list for the names found here and the names returned by enumerate()
	foundNames := make([]string, len(ifis))
	monitorNames := make([]string, len(ifis))

	// Populate the list
	for _, ifi := range ifis {
		foundNames = append(foundNames, ifi.Name)
	}

	// Create a monitor
	monitor := newMonitor(50 * time.Millisecond)
	defer monitor.Stop()

	// Wait for the interfaces to be enumerated
	timeout := time.After(150 * time.Millisecond)

loop:
	for {
		select {
		case name := <-monitor.InterfaceAdded:
			monitorNames = append(monitorNames, name)
		case <-monitor.InterfaceRemoved:
			t.Fatal("Interface removed unexpectedly")
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
