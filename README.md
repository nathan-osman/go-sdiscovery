## go-sdiscovery

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

[TODO]
