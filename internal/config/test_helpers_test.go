package config

import (
	"os"
	"strings"
	"testing"
)

// envSnapshot captures the current process environment so a test can restore
// it during cleanup. Returns a function that re-applies the snapshot.
func envSnapshot(t *testing.T) func() {
	t.Helper()
	saved := append([]string{}, os.Environ()...)
	return func() {
		os.Clearenv()
		for _, kv := range saved {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				continue
			}
			_ = os.Setenv(parts[0], parts[1])
		}
	}
}
