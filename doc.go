// To use go-sdiscovery in your project, include the following import:
//
//     import "github.com/nathan-osman/go-sdiscovery"
//
// Usage
//
// All interaction with the library takes place through an instance of Service,
// which is created in the following manner:
//
//     s := New(ServiceConfig{
//         PollInterval: 1*time.Minute,
//         PingInterval: 2*time.Second,
//         PeerTimeout:  8*time.Second,
//         Port:         1234,
//         ID:           "machine01",
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
// The service can be shutdown by invoking the Stop() method:
//
//     s.Stop()

package sdiscovery
