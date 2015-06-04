package util

import (
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

	// Create a new channel for enumeration.
	enumChan := make(chan time.Time)
	defer close(enumChan)

	// Use a variable to store a map of strings for the EnumFunc.
	strMap := StrMap{"a": nil}

	// Create an enumerator that returns the current value of the map.
	strEnum := NewStrEnum(enumChan, func() (StrMap, error) {
		return strMap, nil
	})

	// Enumerate the map.
	enumChan <- time.Now()

	// Attempt to read the second value.
	if !readStringFromChannel("a", strEnum.StringAdded) {
		t.Fatal("Unable to read initial value from channel")
	}

	// Remove the first value and enumerate again.
	strMap = StrMap{}
	enumChan <- time.Now()

	// Ensure the first value was removed.
	if !readStringFromChannel("a", strEnum.StringRemoved) {
		t.Fatal("Initial value was not removed")
	}
}
