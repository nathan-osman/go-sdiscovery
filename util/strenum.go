package util

import (
	"container/list"
	"log"
	"reflect"
	"time"
)

// Map of strings.
type StrMap map[string]interface{}

// Function that enumerates a list of strings.
type EnumFunc func() (StrMap, error)

// StrEnum provides a simple mechanism for periodically invoking a function
// that enumerates strings and indicating when items are added or removed.
type StrEnum struct {
	StringAdded   chan string // notify when a string is added
	StringRemoved chan string // notify when a string is removed
}

// Create a new enumerator with the specified enumeration function. The
// enumeration process will run each time a value is received from enumChan
// until it is closed. Note that enumFunc must not modify the map it returns.
func NewStrEnum(enumChan <-chan time.Time, enumFunc EnumFunc) *StrEnum {

	// Create a new enumerator
	s := &StrEnum{
		StringAdded:   make(chan string),
		StringRemoved: make(chan string),
	}

	// Launch a separate goroutine to perform the enumeration
	go s.run(enumChan, enumFunc)

	return s
}

// Continually invoke the enumerator until stopped.
func (s *StrEnum) run(enumChan <-chan time.Time, enumFunc EnumFunc) {

	// Map of strings from the previous enumeration and lists of items that
	// need to be sent on one of the notification channels
	oldStrings := StrMap{}
	stringsAdded, stringsRemoved := list.New(), list.New()

	for {

		// Use constants to make it a bit easier to see what's going on.
		const (
			enumCase = iota
			addCase
			removeCase
			numCases
		)

		// A runtime select is needed here in order to avoid a deadlock.
		// Whenever enumeration indicates that notifications should be sent,
		// this needs to take place within the same select as enumChan. Since
		// sending on these channels can block, there needs to be a way to
		// abort the send if the enumChan is closed.
		cases := make([]reflect.SelectCase, numCases)
		cases[enumCase].Dir = reflect.SelectRecv
		cases[enumCase].Chan = reflect.ValueOf(enumChan)
		cases[addCase].Dir = reflect.SelectSend
		cases[removeCase].Dir = reflect.SelectSend

		// If there are values waiting to be sent on either of the two channels
		// then fill in the appropriate fields in each of the cases.
		if stringsAdded.Len() != 0 {
			cases[addCase].Chan = reflect.ValueOf(s.StringAdded)
			cases[addCase].Send = reflect.ValueOf(stringsAdded.Front().Value.(string))
		}
		if stringsRemoved.Len() != 0 {
			cases[removeCase].Chan = reflect.ValueOf(s.StringRemoved)
			cases[removeCase].Send = reflect.ValueOf(stringsRemoved.Front().Value.(string))
		}

		// Perform the select. The value received is ignored.
		i, _, ok := reflect.Select(cases)

		// Perform the appropriate action based on the case index.
		switch i {
		case enumCase:

			// If the receive was successful, perform another enumeration.
			// Otherwise, the channel was closed and the loop should quit.
			if ok {

				// Log and ignore any errors
				if newStrings, err := enumFunc(); err != nil {
					log.Println("[ERR]", err)
				} else {
					stringsAdded.PushBackList(s.compare(oldStrings, newStrings))
					stringsRemoved.PushBackList(s.compare(newStrings, oldStrings))
					oldStrings = newStrings
				}
			} else {
				break
			}

		case addCase:
			stringsAdded.Remove(stringsAdded.Front())
		case removeCase:
			stringsRemoved.Remove(stringsRemoved.Front())
		}
	}

	close(s.StringAdded)
	close(s.StringRemoved)
}

// Compare two maps and return a list of any changes.
func (s *StrEnum) compare(a, b StrMap) *list.List {

	// Create an empty list for new items.
	l := list.New()

	// Check each of the items in map b to see if they exist in map a. If not,
	// add them to the list.
	for item, _ := range b {
		if _, exists := a[item]; !exists {
			l.PushBack(item)
		}
	}

	return l
}
