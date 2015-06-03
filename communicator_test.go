package sdiscovery

/*

// Test that connections are added correctly for interfaces
func Test_addInterface(t *testing.T) {

	// Ensure that at least one valid network interface exists
	ifi, err := findInterfaceWithFlags(net.FlagBroadcast | net.FlagMulticast)
	if err != nil {
		t.Fatal(err)
	}

	// Check to see if one was found
	if ifi == nil {
		t.Skip("No broadcast or multicast interface found")
	}

	// Create a new communicator
	comm := newCommunicator(time.Second, 0)
	defer comm.Stop()

	// Wait a little bit
	<-time.After(50 * time.Millisecond)

	// Ensure that at least one socket was created
	if len(comm.connections) == 0 {
		t.Fatal("No connections were established")
	}
}

*/
