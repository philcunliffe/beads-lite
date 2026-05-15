// Package doltserver is a no-op shim retained so existing call sites can
// compile under the SQLite-only build. The lite build never runs an
// embedded Dolt server.
package doltserver

import "path/filepath"

const (
	// GlobalDatabaseName matches the historical default but has no behavioural
	// meaning in lite mode.
	GlobalDatabaseName = "beads"

	// PIDFileName is retained for callers that report on legacy dirs.
	PIDFileName = "dolt-server.pid"

	// ServerModeExternal is the only mode the lite build acknowledges.
	ServerModeExternal = "external"

	// SharedServerDir is the legacy directory name for shared-server state.
	SharedServerDir = ".beads-server"
)

// Config captures the legacy Dolt server configuration. Fields are kept for
// source compatibility; nothing reads them at runtime.
type Config struct {
	Port int
}

// DefaultConfig returns a zero-value config so callers can read Port without
// panicking.
func DefaultConfig(_ string) Config {
	return Config{}
}

// IsSharedServerMode is always false in lite mode.
func IsSharedServerMode() bool { return false }

// IsRunning is always false in lite mode.
func IsRunning(_ string) bool { return false }

// Start is a no-op in lite mode.
func Start(_ string, _ Config) error { return nil }

// ResolveDoltDir returns the legacy Dolt directory path. The directory may
// not exist in lite installations.
func ResolveDoltDir(beadsDir string) string {
	return filepath.Join(beadsDir, "dolt")
}

// ResolveServerDir returns the legacy shared-server directory path.
func ResolveServerDir(beadsDir string) string {
	return filepath.Join(beadsDir, SharedServerDir)
}

// IsPreV56DoltDir always reports false; lite installations never have
// pre-v56 Dolt directories to recover.
func IsPreV56DoltDir(_ string) bool { return false }

// RecoverPreV56DoltDir is a no-op in lite mode.
func RecoverPreV56DoltDir(_ string) (bool, error) { return false, nil }
