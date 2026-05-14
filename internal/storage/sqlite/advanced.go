//go:build sqlite_lite

package sqlite

import (
	"context"
	"database/sql"

	"github.com/steveyegge/beads/internal/storage/issueops"
	"github.com/steveyegge/beads/internal/types"
)

func (s *Store) CheckEligibility(ctx context.Context, issueID string, tier int) (bool, string, error) {
	var eligible bool
	var reason string
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		eligible, reason, err = issueops.CheckEligibilityInTx(ctx, regularTx, issueID, tier)
		return err
	})
	return eligible, reason, err
}

func (s *Store) ApplyCompaction(ctx context.Context, issueID string, tier int, originalSize int, _ int, commitHash string) error {
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		return issueops.ApplyCompactionInTx(ctx, regularTx, issueID, tier, originalSize, commitHash)
	})
}

func (s *Store) GetTier1Candidates(ctx context.Context) ([]*types.CompactionCandidate, error) {
	var result []*types.CompactionCandidate
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetTier1CandidatesInTx(ctx, regularTx)
		return err
	})
	return result, err
}

func (s *Store) GetTier2Candidates(ctx context.Context) ([]*types.CompactionCandidate, error) {
	var result []*types.CompactionCandidate
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetTier2CandidatesInTx(ctx, regularTx)
		return err
	})
	return result, err
}

func (s *Store) GetRepoMtime(ctx context.Context, repoPath string) (int64, error) {
	var result int64
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetRepoMtimeInTx(ctx, ignoredTx, repoPath)
		return err
	})
	return result, err
}

func (s *Store) SetRepoMtime(ctx context.Context, repoPath, jsonlPath string, mtimeNs int64) error {
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		return issueops.SetRepoMtimeInTx(ctx, ignoredTx, repoPath, jsonlPath, mtimeNs)
	})
}

func (s *Store) ClearRepoMtime(ctx context.Context, repoPath string) error {
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		return issueops.ClearRepoMtimeInTx(ctx, ignoredTx, repoPath)
	})
}

func (s *Store) GetMoleculeLastActivity(ctx context.Context, moleculeID string) (*types.MoleculeLastActivity, error) {
	var result *types.MoleculeLastActivity
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetMoleculeLastActivityInTx(ctx, regularTx, moleculeID)
		return err
	})
	return result, err
}

func (s *Store) GetStaleIssues(ctx context.Context, filter types.StaleFilter) ([]*types.Issue, error) {
	var result []*types.Issue
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetStaleIssuesInTx(ctx, regularTx, filter)
		return err
	})
	return result, err
}
