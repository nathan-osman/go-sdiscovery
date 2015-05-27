## go-sdiscovery

[![MIT License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](http://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/nathan-osman/go-sdiscovery?status.svg)](https://godoc.org/github.com/nathan-osman/go-sdiscovery)
[![Build Status](https://travis-ci.org/nathan-osman/go-sdiscovery.svg)](https://travis-ci.org/nathan-osman/go-sdiscovery)

This library provides an extremely simple API that abstracts the process of registering a service available over the local network and discovering other peers providing the service. This is accomplished by sending broadcast (IPv4) and multicast (IPv6) packets at regular intervals over connected network interfaces.

**Note:** go-sdiscovery does not implement authentication or encryption. Therefore, *it should not be used to transmit sensitive data* and *all data received from other peers should be considered untrusted*. These are both beyond the scope of this library.

### Setup

Use the following command to download the source code for go-sdiscovery:

    go fetch github.com/nathan-osman/go-sdiscovery

To use go-sdiscovery in your project, add the following import:

    import "github.com/nathan-osman/go-sdiscovery"

### Usage

All interaction with the library takes place through a `Service`, which is created in the following manner:

    // Custom data sent to other peers during discovery
    type MyData struct {
        Data string
    }

    // Create the service
    svc := sdiscovery.NewService(
        "machine1",             // unique identifier
        &MyData{"custom data"}, // custom data for this peer
        30 * time.Second,       // timeout before a peer is considered unreachable
    )

At this point, the service will begin sending broadcast packets on all local network interfaces and listening for packets from other peers. The service provides a property that maps each peer's unique identifier to its addresses and user data:

    for id, peer := range svc.Peers {
        fmt.Printf("Peer (%s): %v\n", id, peer)
    }

The service also provides two channels that provide notification when a peer is added or removed:

    for {
        select {
        case <-svc.PeerAdded:
            fmt.Println("New peer added!")
        case <-svc.PeerRemoved:
            fmt.Println("Peer removed!")
        }
    }

Each peer provides a list of its IP addresses, which is guaranteed to contain at least one item. These are suitable for connecting to a service running on the peer:

    peer := svc.Peers["peer_name"]
    conn, err:= net.DialTCP(
        "tcp", nil,
        &net.TCPAddr{
            IP:   peer.Addresses[0],
            Port: 80,
        },
    )

When you are done with the service, it can be shutdown with:

    svc.Shutdown()
