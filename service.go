package sdiscovery

import (
	"time"
)

// ServiceConfig contains the parameters that control how the service behaves.
// Note that it is important to keep the size of UserData to a minimum since
// the entire struct is sent in each packet.
type ServiceConfig struct {
	PollInterval time.Duration // time between polling for network interfaces
	PingInterval time.Duration // time between pings on the network
	PeerTimeout  time.Duration // time after which a peer is considered unreachable
	Port         int           // port used for broadcast and multicast
	ID           string        // unique identifier for the current machine
	UserData     struct{}      // data sent with each packet to other peers
}

// Service provided on the local network
type Service struct {
	PeerAdded   chan Peer
	PeerRemoved chan Peer
	Peers       []Peer
	stop        chan struct{}
	config      *ServiceConfig
}

// Create a new service with the specified peer timeout and user data.
func NewService(config *ServiceConfig) *Service {

	// Create the new service
	s := &Service{
		PeerAdded:   make(chan Peer),
		PeerRemoved: make(chan Peer),
		stop:        make(chan struct{}),
		config:      config,
	}

	// Create a ticker to schedule timeout checks
	ticker := time.NewTicker(s.config.PeerTimeout)

	// Spawn a new goroutine to check for peer timeouts
	go func() {
		for {
			s.checkPeerTimeouts()
			<-ticker.C
		}
	}()

	return s
}

// Check for peer timeouts
func (s *Service) checkPeerTimeouts() {

	// Create an empty slice pointing to the old array and filter the
	// peers based on whether they have expired or not
	peers := s.Peers[:0]
	for _, peer := range s.Peers {

		// Update the peer and check if it has expired
		peer.Update(s.config.PeerTimeout)
		if peer.HasExpired() {
			s.PeerRemoved <- peer
		} else {
			peers = append(peers, peer)
		}
	}
	s.Peers = peers
}
