package doctor

import (
	"context"
	"sync"

	"github.com/steveyegge/beads/internal/storage/dolt"
)

// SharedStore lazily opens the SQLite-backed store for the doctor run and
// keeps a single connection alive across all checks.
type SharedStore struct {
	repoPath string
	beadsDir string

	once   sync.Once
	openErr error
	store  *dolt.DoltStore
}

// NewSharedStore returns a SharedStore rooted at the given repo path.
func NewSharedStore(repoPath string) *SharedStore {
	return &SharedStore{
		repoPath: repoPath,
		beadsDir: ResolveBeadsDirForRepo(repoPath),
	}
}

// Store opens the database on first use and returns the cached store.
func (s *SharedStore) Store() *dolt.DoltStore {
	s.once.Do(func() {
		store, err := dolt.NewFromConfigWithCLIOptions(context.Background(), s.beadsDir, &dolt.Config{ReadOnly: true})
		if err != nil {
			s.openErr = err
			return
		}
		s.store = store
	})
	return s.store
}

// Err returns the deferred open error, if any.
func (s *SharedStore) Err() error {
	s.Store()
	return s.openErr
}

// BeadsDir returns the resolved beads directory.
func (s *SharedStore) BeadsDir() string {
	return s.beadsDir
}

// Close closes the underlying store if one was opened.
func (s *SharedStore) Close() error {
	if s.store == nil {
		return nil
	}
	return s.store.Close()
}

// beadsDirFromSharedStore returns the beads dir tracked by the shared store,
// falling back to per-path resolution when ss is nil.
func beadsDirFromSharedStore(path string, ss *SharedStore) string {
	if ss != nil && ss.beadsDir != "" {
		return ss.beadsDir
	}
	return ResolveBeadsDirForRepo(path)
}

// sharedStoreBeadsDir returns the beads dir that the shared store is rooted at.
func sharedStoreBeadsDir(ss *SharedStore) string {
	if ss == nil {
		return ""
	}
	return ss.beadsDir
}

// sharedStoreNeedsLocalDoltDir is always true in the lite build: the local
// SQLite file under .beads/ is the only thing the doctor needs to find.
func sharedStoreNeedsLocalDoltDir(_ string) bool { return true }
