package sdiscovery

func Test_stuff(t *testing.T) {

}

/*

func Test_checkPeerTimeouts(t *testing.T) {

	// Create a new service with a peer timeout of 50 ms
	svc := NewService(ServiceConfig{
		PeerTimeout: 50 * time.Millisecond,
	})

	// Manually inject the peer into the service
	svc.peers[""] = Peer{}

	// Wait for the peer to be removed
	timeout := time.After(100 * time.Millisecond)

	select {
	case <-svc.PeerRemoved:
	case <-timeout:
		t.Fatal("Timeout reached")
	}

	// Check the number of peers
	if len(svc.peers) != 0 {
		t.Fatal("Peer not removed")
	}
}

*/
