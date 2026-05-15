//go:build sqlite_lite

package main

import (
	"context"
	"path/filepath"

	"github.com/steveyegge/beads/internal/configfile"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/storage/dolt"
	sqlitestore "github.com/steveyegge/beads/internal/storage/sqlite"
)

const sqliteLiteDatabaseFile = "beads.sqlite3"

func usesSQLServer() bool   { return false }
func usesProxiedServer() bool { return false }

func newDoltStore(ctx context.Context, cfg *dolt.Config) (storage.DoltStorage, error) {
	beadsDir := cfg.BeadsDir
	if beadsDir == "" && cfg.Path != "" {
		beadsDir = filepath.Dir(cfg.Path)
	}
	return openSQLiteLiteStore(ctx, beadsDir)
}

func newDoltStoreFromConfig(ctx context.Context, beadsDir string) (storage.DoltStorage, error) {
	return openSQLiteLiteStore(ctx, beadsDir)
}

func newReadOnlyStoreFromConfig(ctx context.Context, beadsDir string) (storage.DoltStorage, error) {
	return openSQLiteLiteStore(ctx, beadsDir)
}

func openSQLiteLiteStore(ctx context.Context, beadsDir string) (storage.DoltStorage, error) {
	cfg, err := configfile.Load(beadsDir)
	if err != nil {
		return nil, err
	}
	dbName := sqliteLiteDatabaseFile
	if cfg != nil && cfg.Database != "" {
		dbName = cfg.Database
	}
	if filepath.IsAbs(dbName) {
		return sqlitestore.Open(ctx, dbName)
	}
	return sqlitestore.Open(ctx, filepath.Join(beadsDir, dbName))
}
