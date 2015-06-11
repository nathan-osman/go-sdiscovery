package ifienum

import (
	"testing"
	"time"
)

// Ensure that the enumerator can be created and destroyed without issue.
func Test_New(t *testing.T) {
	ticker := time.NewTicker(time.Second)
	New(ticker.C)
	ticker.Stop()
}
