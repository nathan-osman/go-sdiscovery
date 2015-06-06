package comm

import (
	"testing"
	"time"
)

// Ensure that the Communicator class can be instantiated and terminated.
func Test_Communicator(t *testing.T) {
	c := NewCommunicator(time.Second, 8000)
	defer c.Stop()
}
