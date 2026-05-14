//go:build sqlite_lite

package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/steveyegge/beads/internal/storage/issueops"
	"github.com/steveyegge/beads/internal/types"
)

func (s *Store) GetStatistics(ctx context.Context) (*types.Statistics, error) {
	stats := &types.Statistics{}
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		if err := issueops.ScanIssueCountsInTx(ctx, regularTx, stats); err != nil {
			return err
		}

		blockedIDs, _, err := issueops.ComputeBlockedIDsInTx(ctx, regularTx, true)
		if err != nil {
			return err
		}
		stats.BlockedIssues = len(blockedIDs)
		stats.ReadyIssues = stats.OpenIssues - stats.BlockedIssues
		if stats.ReadyIssues < 0 {
			stats.ReadyIssues = 0
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("sqlite: get statistics: %w", err)
	}
	return stats, nil
}
