package sdiscovery

import (
	"time"

	"github.com/nathan-osman/go-sdiscovery/conn"
	"github.com/nathan-osman/go-sdiscovery/peer"
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
	UserData     []byte        // data sent with each packet to other peers
}

// Service sends and receives packets on local network interfaces in order to
// discover other peers providing the service and announce its presence.
type Service struct {
	PeerAdded   chan string // indicates that a new peer was found
	PeerRemoved chan string // indicates that an existing peer has timed out
	stop        chan interface{}
	peers       map[string]*peer.Peer
	config      ServiceConfig
}

// NewService creates a new Service instance with the specified configuration.
func NewService(config ServiceConfig) *Service {

	s := &Service{
		PeerAdded:   make(chan string),
		PeerRemoved: make(chan string),
		stop:        make(chan interface{}),
		peers:       make(map[string]*peer.Peer),
		config:      config,
	}

	// Spawn a new goroutine for:
	// - sending and receiving packets
	// - checking for expired peers
	go s.run()

	return s
}

// Process pings and expire peers.
func (s *Service) run() {

	// Create a communicator for sending and receiving packets.
	communicator := conn.NewCommunicator(s.config.PollInterval, s.config.Port)
	defer communicator.Stop()

	// Create a ticker for sending pings.
	pingTicker := time.NewTicker(s.config.PingInterval)
	defer pingTicker.Stop()

	// Create a ticker for timeout checks.
	timeoutTicker := time.NewTicker(s.config.PeerTimeout)
	defer timeoutTicker.Stop()

	pkt := &conn.Packet{
		ID:       s.config.ID,
		UserData: s.config.UserData,
	}

	for {

		select {
		case p := <-communicator.PacketChan:
			s.processPacket(p)
		case <-pingTicker.C:
			communicator.Send(pkt)
		case <-timeoutTicker.C:
			s.processPeers()
		case <-s.stop:
			return
		}
	}
}

// Process a packet received.
func (s *Service) processPacket(pkt *conn.Packet) {

	// Only process packets that did not originate from here.
	if pkt.ID != s.config.ID {

		// If the peer does not exist, create a new one.
		_, exists := s.peers[pkt.ID]
		if !exists {
			s.peers[pkt.ID] = &peer.Peer{}
		}

		// Update the peer.
		s.peers[pkt.ID].Ping(pkt, time.Now())

		// Write to the PeerAdded channel if the peer is new.
		if !exists {
			s.PeerAdded <- pkt.ID
		}
	}
}

// Check for peer timeouts.
func (s *Service) processPeers() {

	curTime := time.Now()

	for id, peer := range s.peers {
		if peer.IsExpired(s.config.PeerTimeout, curTime) {

			// Write to the PeerRemoved channel and delete it.
			s.PeerRemoved <- id
			delete(s.peers, id)
		}
	}
}

// Stop the service.
func (s *Service) Stop() {
	close(s.stop)
}
