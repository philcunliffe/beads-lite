//go:build sqlite_lite

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/steveyegge/beads/internal/config"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"

	_ "modernc.org/sqlite"
)

var _ storage.DoltStorage = (*Store)(nil)
var _ storage.RawDBAccessor = (*Store)(nil)
var _ storage.StoreLocator = (*Store)(nil)
var _ storage.LifecycleManager = (*Store)(nil)

// Store implements the beads storage interfaces on top of one local SQLite DB.
type Store struct {
	db       *sql.DB
	path     string
	beadsDir string
	closed   atomic.Bool
}

var errClosed = errors.New("sqlite: store is closed")

// Open opens or creates a SQLite beads database at path.
func Open(ctx context.Context, path string) (*Store, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("sqlite: resolving path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(absPath), config.BeadsDirPerm); err != nil {
		return nil, fmt.Errorf("sqlite: creating data directory: %w", err)
	}

	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	s := &Store{
		db:       db,
		path:     absPath,
		beadsDir: filepath.Dir(absPath),
	}
	if err := s.configure(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := s.initSchema(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) configure(ctx context.Context) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA busy_timeout = 5000",
	}
	for _, q := range pragmas {
		if _, err := s.db.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("sqlite: %s: %w", q, err)
		}
	}
	return nil
}

func (s *Store) withConn(ctx context.Context, commit bool, fn func(regularTx, ignoredTx *sql.Tx) error) (err error) {
	if s.closed.Load() {
		return errClosed
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sqlite: begin tx: %w", err)
	}
	if fnErr := fn(tx, tx); fnErr != nil {
		return errors.Join(fnErr, tx.Rollback())
	}
	if !commit {
		return tx.Rollback()
	}
	return tx.Commit()
}

func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) UnderlyingDB() *sql.DB {
	return s.db
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) CLIDir() string {
	return s.beadsDir
}

func (s *Store) IsClosed() bool {
	return s.closed.Load()
}

func (s *Store) Close() error {
	if s.closed.Swap(true) {
		return nil
	}
	return s.db.Close()
}

func unsupported(name string) error {
	return fmt.Errorf("%w: sqlite backend does not support %s", storage.ErrUnsupportedCapability, name)
}

// Version/history/remote/sync APIs are intentionally unsupported by bd-lite.
func (s *Store) Branch(ctx context.Context, name string) error     { return unsupported("branches") }
func (s *Store) Checkout(ctx context.Context, branch string) error { return unsupported("checkout") }
func (s *Store) CurrentBranch(ctx context.Context) (string, error) { return "main", nil }
func (s *Store) DeleteBranch(ctx context.Context, branch string) error {
	return unsupported("branches")
}
func (s *Store) ListBranches(ctx context.Context) ([]string, error)            { return []string{"main"}, nil }
func (s *Store) Commit(ctx context.Context, message string) error              { return nil }
func (s *Store) CommitWithConfig(ctx context.Context, message string) error    { return nil }
func (s *Store) CommitPending(ctx context.Context, actor string) (bool, error) { return false, nil }
func (s *Store) CommitExists(ctx context.Context, commitHash string) (bool, error) {
	return false, unsupported("commit lookup")
}
func (s *Store) GetCurrentCommit(ctx context.Context) (string, error) { return "sqlite", nil }
func (s *Store) Status(ctx context.Context) (*storage.Status, error)  { return &storage.Status{}, nil }
func (s *Store) Log(ctx context.Context, limit int) ([]storage.CommitInfo, error) {
	return nil, unsupported("history")
}
func (s *Store) Merge(ctx context.Context, branch string) ([]storage.Conflict, error) {
	return nil, unsupported("merge")
}
func (s *Store) GetConflicts(ctx context.Context) ([]storage.Conflict, error) {
	return nil, unsupported("conflicts")
}
func (s *Store) ResolveConflicts(ctx context.Context, table, strategy string) error {
	return unsupported("conflicts")
}
func (s *Store) History(ctx context.Context, issueID string) ([]*storage.HistoryEntry, error) {
	return nil, unsupported("history")
}
func (s *Store) AsOf(ctx context.Context, issueID string, ref string) (*types.Issue, error) {
	return nil, unsupported("as-of")
}
func (s *Store) Diff(ctx context.Context, fromRef, toRef string) ([]*storage.DiffEntry, error) {
	return nil, unsupported("diff")
}
func (s *Store) AddRemote(ctx context.Context, name, url string) error         { return unsupported("remotes") }
func (s *Store) RemoveRemote(ctx context.Context, name string) error           { return unsupported("remotes") }
func (s *Store) HasRemote(ctx context.Context, name string) (bool, error)      { return false, nil }
func (s *Store) ListRemotes(ctx context.Context) ([]storage.RemoteInfo, error) { return nil, nil }
func (s *Store) Push(ctx context.Context) error                                { return unsupported("push") }
func (s *Store) Pull(ctx context.Context) error                                { return unsupported("pull") }
func (s *Store) ForcePush(ctx context.Context) error                           { return unsupported("push") }
func (s *Store) PushRemote(ctx context.Context, remote string, force bool) error {
	return unsupported("push")
}
func (s *Store) PullRemote(ctx context.Context, remote string) error { return unsupported("pull") }
func (s *Store) Fetch(ctx context.Context, peer string) error        { return unsupported("fetch") }
func (s *Store) PushTo(ctx context.Context, peer string) error       { return unsupported("push") }
func (s *Store) PullFrom(ctx context.Context, peer string) ([]storage.Conflict, error) {
	return nil, unsupported("pull")
}
func (s *Store) Sync(ctx context.Context, peer string, strategy string) (*storage.SyncResult, error) {
	return nil, unsupported("sync")
}
func (s *Store) SyncStatus(ctx context.Context, peer string) (*storage.SyncStatus, error) {
	return nil, unsupported("sync")
}
func (s *Store) AddFederationPeer(ctx context.Context, peer *storage.FederationPeer) error {
	return unsupported("federation")
}
func (s *Store) GetFederationPeer(ctx context.Context, name string) (*storage.FederationPeer, error) {
	return nil, unsupported("federation")
}
func (s *Store) ListFederationPeers(ctx context.Context) ([]*storage.FederationPeer, error) {
	return nil, unsupported("federation")
}
func (s *Store) RemoveFederationPeer(ctx context.Context, name string) error {
	return unsupported("federation")
}
