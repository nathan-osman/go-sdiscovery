package sdiscovery

import (
	"strconv"
	"testing"
	"time"
)

// Attempt to read the specified string from the channel.
func readStringFromChannel(orig string, c chan string) bool {
	select {
	case v := <-c:
		return v == orig
	case <-time.After(1 * time.Second):
		return false
	}
}

// Ensure that notifications are properly dispatched as strings are added and
// removed from a map.
func Test_StrEnum(t *testing.T) {

	// In order to ensure different items are returned each time the enumerate
	// function is invoked, use a local integer variable
	i := 0

	// Create a new StrEnum with the specified interval and supply a function
	// that returns a predefined sequence of maps
	c := make(chan time.Time)
	s := NewStrEnum(c, func() StringMap {
		i++
		return StringMap{strconv.FormatInt(int64(i), 10): nil}
	})

	// Attempt to read the first value
	c <- time.Now()
	if !readStringFromChannel("2", s.StringAdded) ||
		!readStringFromChannel("1", s.StringRemoved) {
		t.Fatal("Unable to read expected values from channels")
	}
}
