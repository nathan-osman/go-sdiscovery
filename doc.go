// All interaction with the library takes place through an instance of Service,
// which is created in the following manner:
//
//     s := New(ServiceConfig{
//         PollInterval: 1*time.Minute,
//         PingInterval: 2*time.Second,
//         PeerTimeout:  8*time.Second,
//         Port:         1234,
//         ID:           "machine01",
//         UserData:     []byte("data"),
//     })
//
// At this point, the service will begin sending broadcast and multicast
// packets on all appropriate network interfaces and listening for packets from
// other peers. The service provides two channels that provide notifications
// when peers are added or removed:
//
//     for {
//         select {
//         case id := <- s.PeerAdded:
//             fmt.Printf("Peer %s added!\n", id)
//         case id := <- s.PeerRemoved:
//             fmt.Printf("Peer %s removed!\n", id)
//         }
//     }
//
// Once you have a peer ID, you can use it to retrieve the custom user data for
// that specific peer:
//
//     data, _ := s.PeerUserData(id)
//     fmt.Printf("UserData: %s\n", data)
//
// If you need to connect to the peer, it is possible to obtain a slice of IP
// addresses for the peer. As packets are received from the peer, the IP
// address and timestamp are recored. This allows the service to determine
// the best IP address for contacting the peer.
//
//     addrs, _ := s.PeerAddrs(id)
//     for _, a := range addrs {
//         fmt.Printf("- %s", a)
//     }
//
// Note that you may want to filter the addresses since the slice may contain
// both IPv4 and IPv6 addresses.
//
// The service can be shutdown by invoking the Stop() method:
//
//     s.Stop()
//
package sdiscovery
