package conn

/*

// Attempt to find an interface with the specified flags
func findInterfaceWithFlags(flags net.Flags) (*net.Interface, error) {

	// Obtain the list of interfaces
	ifis, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Return the first one that matches
	for _, ifi := range ifis {
		if ifi.Flags&flags != 0 {
			return &ifi, nil
		}
	}

	// None matched - return nil
	return nil, nil
}

// Send and receive a packet
func sendAndReceivePacket(ifi *net.Interface, multicast bool) error {

	packetReceived := make(chan packet)
	var waitGroup sync.WaitGroup

	// Create the connection with a randomly chosen port
	conn, err := newConnection(packetReceived, &waitGroup, multicast, ifi, 0)
	if err != nil {
		return err
	}

	defer conn.Stop()

	// Send a packet
	packet := []byte(`test`)
	if err := conn.Send(packet); err != nil {
		return err
	}

	// Receive the packet
	select {
	case b := <-packetReceived:
		if !bytes.Equal(b.Data, packet) {
			return errors.New("Packet contents do not match")
		}
	case <-time.NewTicker(50 * time.Millisecond).C:
		return errors.New("Timeout waiting for broadcast packet")
	}

	return nil
}

// Test that packets are correctly sent and received via broadcast
func Test_connection_broadcast(t *testing.T) {

	// Attempt to find a broadcast interface
	ifi, err := findInterfaceWithFlags(net.FlagBroadcast)
	if err != nil {
		t.Fatal(err)
	}

	// Skip the test if none was found
	if ifi == nil {
		t.Skip("No broadcast interface found")
	}

	// Run the test
	if err := sendAndReceivePacket(ifi, false); err != nil {
		t.Fatal(err)
	}
}

// Test that packets are correctly sent and received via multicast
func Test_connection_multicast(t *testing.T) {

	// Attempt to find a multicast interface
	ifi, err := findInterfaceWithFlags(net.FlagMulticast)
	if err != nil {
		t.Fatal(err)
	}

	// Skip the test if none was found
	if ifi == nil {
		t.Skip("No multicast interface found")
	}

	// Run the test
	if err := sendAndReceivePacket(ifi, true); err != nil {
		t.Fatal(err)
	}
}

*/
