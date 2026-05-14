//go:build sqlite_lite

package sqlite

import (
	"context"
	"database/sql"

	"github.com/steveyegge/beads/internal/storage/issueops"
)

func (s *Store) GetNextChildID(ctx context.Context, parentID string) (string, error) {
	var childID string
	err := s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		childID, err = issueops.GetNextChildIDTx(ctx, regularTx, parentID)
		return err
	})
	return childID, err
}
