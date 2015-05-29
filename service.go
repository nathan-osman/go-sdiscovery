package sdiscovery

import (
	"time"
)

// ServiceConfig contains the parameters that control how the service behaves.
// Note that it is important to keep the size of UserData to a minimum since
// the entire struct is sent in each packet. This struct should not be modified
// after being passed to NewService.
type ServiceConfig struct {
	PollInterval time.Duration // time between polling for network interfaces
	PingInterval time.Duration // time between pings on the network
	PeerTimeout  time.Duration // time after which a peer is considered unreachable
	Port         int           // port used for broadcast and multicast
	ID           string        // unique identifier for the current machine
	UserData     struct{}      // data sent with each packet to other peers
}

// Service sends and receives packets on local network interfaces in order to
// discover other peers providing the service and announce its presence.
type Service struct {
	PeerAdded   chan Peer // indicates that a new peer was found
	PeerRemoved chan Peer // indicates that an existing peer has timed out
	stop        chan struct{}
	peers       map[string]Peer
	config      ServiceConfig
}

// NewService creates a new Service instance with the specified configuration.
func NewService(config ServiceConfig) *Service {

	s := &Service{
		PeerAdded:   make(chan Peer),
		PeerRemoved: make(chan Peer),
		stop:        make(chan struct{}),
		peers:       make(map[string]Peer),
		config:      config,
	}

	// Spawn a new goroutine for:
	// - sending and receiving packets
	// - checking for expired peers
	go s.run()

	return s
}

// Process pings and expire peers
func (s *Service) run() {

	// Create a communicator for sending and receiving packets
	communicator := newCommunicator(s.config.PollInterval, s.config.Port)
	defer communicator.Stop()

	// Create a ticker for sending pings
	pingTicker := time.NewTicker(s.config.PingInterval)
	defer pingTicker.Stop()

	// Create a ticker for timeout checks
	timeoutTicker := time.NewTicker(s.config.PeerTimeout)
	defer timeoutTicker.Stop()

	for {
		select {
		case p := <-communicator.PacketReceived:
			s.processPacket(p)
		case <-pingTicker.C:
			communicator.Send(s.config.ID, s.config.UserData)
		case <-timeoutTicker.C:
			s.processPeers()
		case <-s.stop:
			return
		}
	}
}

// Process a packet received
func (s *Service) processPacket(p packet) {

	// Only process packets that did not originate from here
	if p.ID != s.config.ID {

		// If the peer does not exist, create a new one
		_, exists := s.peers[p.ID]
		if !exists {
			s.peers[p.ID] = Peer{}
		}

		// Update the peer
		s.peers[p.ID].Update(s.config.PeerTimeout)

		// Write to the PeerAdded channel if the peer is new
		if !exists {
			s.PeerAdded <- s.peers[p.ID]
		}
	}
}

// Check for peer timeouts
func (s *Service) processPeers() {
	for id, peer := range s.peers {
		if peer.HasExpired() {

			// Write to the PeerRemoved channel and delete it
			s.PeerRemoved <- peer
			delete(s.peers, id)
		}
	}
}

// Stop the service
func (s *Service) Stop() {
	close(s.stop)
}
