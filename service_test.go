package sdiscovery

import (
	"testing"
	"time"
)

// Ensure that the Service class can be instantiated and terminated.
func Test_Service(t *testing.T) {
	s := New(ServiceConfig{
		PollInterval: time.Second,
		PingInterval: time.Second,
		PeerTimeout:  time.Second,
		Port:         8000,
		ID:           "1234",
		UserData:     nil,
	})
	defer s.Stop()
}
