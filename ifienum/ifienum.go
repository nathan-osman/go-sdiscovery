package ifienum

import (
	"container/list"
	"log"
	"net"
	"reflect"
	"time"
)

// Constants for the dynamic select in run().
const (
	enumCase = iota
	addCase
	removeCase
	numCases // number of select cases
)

// Use a map for storing the list of interface names to make set comparison
// based on the keys possible.
type ifiMap map[string]struct{}

// IfiEnum continously retrieves the list of attached network interfaces and
// writes to the IfiAddedChan and IfiRemovedChan channels when interfaces are
// added and removed. Errors are currently logged but ignored.
type IfiEnum struct {
	IfiAddedChan   chan string
	IfiRemovedChan chan string
	enumChan       <-chan time.Time
}

// Create a new enumerator that polls for new interfaces whenever a value is
// received on the provided channel. Enumeration will continue until the
// channel is closed.
func New(enumChan <-chan time.Time) *IfiEnum {
	i := &IfiEnum{
		IfiAddedChan:   make(chan string),
		IfiRemovedChan: make(chan string),
		enumChan:       enumChan,
	}
	go i.run()
	return i
}

// Continue to enumerate interfaces until the channel is closed. In order to
// avoid blocking when sending on the IfiAdded/Removed channels, a runtime
// select is used and the interface names are buffered in a list.
func (i IfiEnum) run() {
	var (
		ifisAdded   = list.New()
		ifisRemoved = list.New()
	)
	oldIfis := i.enumIfis(ifiMap{})
	ifisAdded.PushBackList(i.compare(ifiMap{}, oldIfis))
loop:
	for {
		cases := i.createCases(ifisAdded, ifisRemoved)
		idx, _, ok := reflect.Select(cases)

		switch idx {
		case enumCase:
			if ok {
				newIfis := i.enumIfis(oldIfis)
				ifisAdded.PushBackList(i.compare(oldIfis, newIfis))
				ifisRemoved.PushBackList(i.compare(newIfis, oldIfis))
				oldIfis = newIfis
			} else {
				break loop
			}
		case addCase:
			ifisAdded.Remove(ifisAdded.Front())
		case removeCase:
			ifisRemoved.Remove(ifisRemoved.Front())
		}
	}
	close(i.IfiAddedChan)
	close(i.IfiRemovedChan)
}

// Create a slice of select cases consisting of the enumeration channel and any
// values waiting to be written to the IfiAdded/Removed channels.
func (i IfiEnum) createCases(ifisAdded, ifisRemoved *list.List) []reflect.SelectCase {
	cases := make([]reflect.SelectCase, numCases)
	cases[enumCase].Dir = reflect.SelectRecv
	cases[enumCase].Chan = reflect.ValueOf(i.enumChan)
	cases[addCase].Dir = reflect.SelectSend
	cases[removeCase].Dir = reflect.SelectSend

	if ifisAdded.Len() != 0 {
		cases[addCase].Chan = reflect.ValueOf(i.IfiAddedChan)
		cases[addCase].Send = reflect.ValueOf(ifisAdded.Front().Value.(string))
	}
	if ifisRemoved.Len() != 0 {
		cases[removeCase].Chan = reflect.ValueOf(i.IfiRemovedChan)
		cases[removeCase].Send = reflect.ValueOf(ifisRemoved.Front().Value.(string))
	}
	return cases
}

// Obtain a list of all current interface names and add them to a map. If an
// error occurs, then the error is logged and the old list of interfaces is
// returned instead.
func (i IfiEnum) enumIfis(oldIfis ifiMap) ifiMap {
	ifis, err := net.Interfaces()
	if err != nil {
		log.Println("[ERR]", err)
		return oldIfis
	}
	newIfis := make(ifiMap)
	for _, ifi := range ifis {
		newIfis[ifi.Name] = struct{}{}
	}
	return newIfis
}

// Compare two maps of interface names and return a list of items that are
// present in map b but not in map a.
func (i IfiEnum) compare(a, b ifiMap) *list.List {
	newNames := list.New()
	for name := range b {
		if _, exists := a[name]; !exists {
			newNames.PushBack(name)
		}
	}
	return newNames
}
