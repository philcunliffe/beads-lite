// Package doltutil is a no-op shim retained so existing callers compile
// under the SQLite-only build. None of these utilities perform any real
// work in lite mode.
package doltutil

import "fmt"

// ServerDSN matches the legacy DSN builder fields.
type ServerDSN struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// String renders a minimal DSN; lite mode never connects to Dolt.
func (d ServerDSN) String() string {
	return fmt.Sprintf("dolt://%s:%d/%s", d.Host, d.Port, d.Database)
}

// Remote captures the historical (name, url) pair for Dolt remotes.
type Remote struct {
	Name string
	URL  string
}

// ListCLIRemotes returns an empty slice in lite mode.
func ListCLIRemotes(_ string) ([]Remote, error) { return nil, nil }

// FindCLIRemote always returns ("", false) in lite mode.
func FindCLIRemote(_ string, _ string) (string, bool, error) { return "", false, nil }

// AddCLIRemote is a no-op in lite mode.
func AddCLIRemote(_ string, _ string, _ string) error { return nil }

// RemoveCLIRemote is a no-op in lite mode.
func RemoveCLIRemote(_ string, _ string) error { return nil }

// ToRemoteNameMap converts a list of remotes into a name -> URL map.
func ToRemoteNameMap(remotes []Remote) map[string]string {
	out := make(map[string]string, len(remotes))
	for _, r := range remotes {
		out[r.Name] = r.URL
	}
	return out
}

// IsSSHURL is a best-effort detection retained for the few callers that
// inspect URL shapes for diagnostics.
func IsSSHURL(url string) bool {
	if len(url) > 4 && url[:4] == "git@" {
		return true
	}
	if len(url) > 6 && url[:6] == "ssh://" {
		return true
	}
	return false
}
