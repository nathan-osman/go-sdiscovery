package sdiscovery

import (
	"errors"
	"sync"
	"time"

	"github.com/nathan-osman/go-sdiscovery/conn"
	"github.com/nathan-osman/go-sdiscovery/peer"
)

type peerMap map[string]*peer.Peer

// ServiceConfig contains the parameters that control how the service behaves.
// Note that it is important to keep the size of UserData to a minimum since
// the entire struct is sent in each packet. Any modifications to this struct
// after passing it to New() will be ignored.
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
	sync.Mutex
	PeerAdded   chan string // indicates that a new peer was found
	PeerRemoved chan string // indicates that an existing peer has timed out
	stopChan    chan interface{}
	peers       peerMap
	config      ServiceConfig
}

// Create a new Service instance with the specified configuration.
func New(config ServiceConfig) *Service {

	s := &Service{
		PeerAdded:   make(chan string),
		PeerRemoved: make(chan string),
		stopChan:    make(chan interface{}),
		peers:       make(peerMap),
		config:      config,
	}

	// Spawn a new goroutine for managing peers.
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
	peerTicker := time.NewTicker(s.config.PeerTimeout)
	defer peerTicker.Stop()

	// Create the packet that will be sent to all peers.
	pkt := &conn.Packet{
		ID:       s.config.ID,
		UserData: s.config.UserData,
	}

	// Continue processing events until explicitly stopped.
	for {
		select {
		case p := <-communicator.PacketChan:
			s.processPacket(p)
		case <-pingTicker.C:
			communicator.Send(pkt)
		case <-peerTicker.C:
			s.processPeers()
		case <-s.stopChan:
			return
		}
	}
}

// Process a packet received from one of the connections.
func (s *Service) processPacket(pkt *conn.Packet) {

	// Obtain exclusive access to the map.
	s.Lock()
	defer s.Unlock()

	// Check the ID on the packet to ensure it does not match this peer.
	if pkt.ID != s.config.ID {

		// If the peer ID is not in the map, then create a new one.
		_, exists := s.peers[pkt.ID]
		if !exists {
			s.peers[pkt.ID] = &peer.Peer{}
		}

		// Update the peer with the packet that was received.
		s.peers[pkt.ID].Ping(pkt, time.Now())

		// If the peer didn't exist in the map prior to this packet, then send
		// the peer ID over the PeerAdded channel.
		if !exists {
			s.PeerAdded <- pkt.ID
		}
	}
}

// Check each of the peers in order to determine if any expired.
func (s *Service) processPeers() {

	// Obtain exclusive access to the map.
	s.Lock()
	defer s.Unlock()

	// Avoid repeated calls to time.Now() by invoking it once here.
	curTime := time.Now()

	for id, peer := range s.peers {
		if peer.IsExpired(s.config.PeerTimeout, curTime) {

			// Send the peer ID over the PeerRemoved channel and remove it.
			s.PeerRemoved <- id
			delete(s.peers, id)
		}
	}
}

// Obtain the custom user data provided by the specified peer.
func (s *Service) PeerUserData(id string) ([]byte, error) {

	// Obtain exclusive access to the map.
	s.Lock()
	defer s.Unlock()

	// Attempt to retrieve the peer from the map.
	p, exists := s.peers[id]
	if !exists {
		return nil, errors.New("Peer does not exist")
	}

	return p.UserData, nil
}

// Stop the service. No more packets will be sent or received and all
// connections will be closed.
func (s *Service) Stop() {
	close(s.stopChan)
}
