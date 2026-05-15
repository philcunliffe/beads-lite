// Package dolt is a SQLite-only shim retained so the rest of the codebase
// can continue to reference the historical "dolt" surface without
// reintroducing the embedded Dolt engine. All operations are backed by
// the local SQLite store.
package dolt

import (
	"context"
	"path/filepath"

	"github.com/steveyegge/beads/internal/configfile"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/storage/sqlite"
)

// Config carries the options the legacy Dolt factory accepted. Most fields
// are accepted for source compatibility and ignored at runtime.
type Config struct {
	Path            string
	BeadsDir        string
	Database        string
	CreateIfMissing bool
	ReadOnly        bool
	ServerMode      bool
	ProxiedServer   bool
	ServerHost      string
	ServerPort      int
	ServerSocket    string
	ServerUser      string
	ServerPassword  string
	ServerTLS       bool
	SyncRemote      string
}

// DoltStore is the SQLite store. The alias preserves the historical name.
type DoltStore = sqlite.Store

// New opens (or creates) the SQLite store described by cfg.
func New(ctx context.Context, cfg *Config) (*DoltStore, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	path := cfg.Path
	if path == "" {
		dir := cfg.BeadsDir
		if dir == "" {
			dir = "."
		}
		dbName := cfg.Database
		if dbName == "" {
			dbName = "beads.sqlite3"
		}
		path = filepath.Join(dir, dbName)
	}
	return sqlite.Open(ctx, path)
}

// NewFromConfig opens the store recorded in beadsDir/metadata.json.
func NewFromConfig(ctx context.Context, beadsDir string) (*DoltStore, error) {
	return openFromBeadsDir(ctx, beadsDir)
}

// NewFromConfigWithOptions matches the old signature; options are honoured
// only insofar as they apply to SQLite (effectively just path resolution).
func NewFromConfigWithOptions(ctx context.Context, beadsDir string, _ *Config) (*DoltStore, error) {
	return openFromBeadsDir(ctx, beadsDir)
}

// NewFromConfigWithCLIOptions is an alias kept for compatibility.
func NewFromConfigWithCLIOptions(ctx context.Context, beadsDir string, _ *Config) (*DoltStore, error) {
	return openFromBeadsDir(ctx, beadsDir)
}

// CleanStaleCircuitBreakerFiles is a no-op in lite mode.
func CleanStaleCircuitBreakerFiles() {}

// ApplyCLIAutoStart is a no-op in lite mode.
func ApplyCLIAutoStart(_ string, _ *Config) {}

// DefaultInfraTypes returns the default infra types from the storage package.
func DefaultInfraTypes() []string {
	return storage.DefaultInfraTypes()
}

// GetBackendFromConfig reports the backend recorded in metadata.json.
func GetBackendFromConfig(beadsDir string) string {
	cfg, err := configfile.Load(beadsDir)
	if err != nil || cfg == nil {
		return ""
	}
	return string(cfg.Backend)
}

func openFromBeadsDir(ctx context.Context, beadsDir string) (*DoltStore, error) {
	cfg, _ := configfile.Load(beadsDir)
	dbName := "beads.sqlite3"
	if cfg != nil && cfg.Database != "" {
		dbName = cfg.Database
	}
	if filepath.IsAbs(dbName) {
		return sqlite.Open(ctx, dbName)
	}
	return sqlite.Open(ctx, filepath.Join(beadsDir, dbName))
}
