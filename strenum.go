package sdiscovery

import (
	"time"
)

// Map of strings.
type StringMap map[string]interface{}

// Function that enumerates a list of strings.
type EnumerateFunc func() StringMap

// StrEnum provides a simple mechanism for periodically invoking a function
// that enumerates strings and indicating when items are added or removed.
type StrEnum struct {
	StringAdded   chan string // notify when a string is added
	StringRemoved chan string // notify when a string is removed
	stopChan      chan interface{}
}

// Create a new enumerator with the specified enumeration function. The
// enumeration process will begin immediately and continue at the specified
// interval.
func NewStrEnum(duration time.Duration, enumFunc EnumerateFunc) *StrEnum {

	// Create a new enumerator
	s := &StrEnum{
		StringAdded:   make(chan string),
		StringRemoved: make(chan string),
		stopChan:      make(chan interface{}),
	}

	// Launch a separate goroutine to perform the enumeration
	go s.run(duration, enumFunc)

	return s
}

// Continually invoke the enumerator until stopped.
func (s *StrEnum) run(duration time.Duration, enumFunc EnumerateFunc) {

	// Map of strings from the previous enumeration
	oldStrings := enumFunc()

	for {
		select {
		case <-time.After(duration):
			newStrings := enumFunc()
			s.compare(oldStrings, newStrings, s.StringAdded)
			s.compare(newStrings, oldStrings, s.StringRemoved)
			oldStrings = newStrings
		case <-s.stopChan:
			close(s.StringAdded)
			close(s.StringRemoved)
			return
		}
	}
}

// Compare two maps and notify of any changes on the specified channel
func (s *StrEnum) compare(a, b map[string]interface{}, notifyChan chan<- string) {
	for item, _ := range b {
		if _, exists := a[item]; !exists {
			select {
			case notifyChan <- item:
			case <-s.stopChan:
			}
		}
	}
}

// Immediately stop the enumerator.
func (s *StrEnum) Stop() {
	s.stopChan <- nil
}
