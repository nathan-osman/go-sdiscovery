package sdiscovery

import (
	"time"
)

// Service available on the local network
type Service struct {
	PeerAdded   chan Peer
	PeerRemoved chan Peer
	Peers       []Peer
	peerTimeout time.Duration
	userData    struct{}
}

// Create a new service with the specified peer timeout and user data.
func NewService(peerTimeout time.Duration, userData struct{}) *Service {

	// Create the new service
	s := &Service{
		PeerAdded:   make(chan Peer),
		PeerRemoved: make(chan Peer),
		peerTimeout: peerTimeout,
		userData:    userData,
	}

	// Create a ticker to schedule timeout checks
	ticker := time.NewTicker(peerTimeout)

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
		peer.Update(s.peerTimeout)
		if peer.HasExpired() {
			s.PeerRemoved <- peer
		} else {
			peers = append(peers, peer)
		}
	}
	s.Peers = peers
}
