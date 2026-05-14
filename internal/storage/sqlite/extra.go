//go:build sqlite_lite

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/storage/issueops"
	"github.com/steveyegge/beads/internal/types"
	"github.com/steveyegge/beads/internal/utils"
)

func (s *Store) GetIssueByExternalRef(ctx context.Context, externalRef string) (*types.Issue, error) {
	var id string
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		id, err = issueops.GetIssueByExternalRefInTx(ctx, regularTx, externalRef)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetIssue(ctx, id)
}

func (s *Store) DeleteIssue(ctx context.Context, id string) error {
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		return issueops.DeleteIssueInTx(ctx, regularTx, ignoredTx, id)
	})
}

func (s *Store) GetDependencies(ctx context.Context, issueID string) ([]*types.Issue, error) {
	var result []*types.Issue
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetDependenciesInTx(ctx, regularTx, issueID)
		return err
	})
	return result, err
}

func (s *Store) GetDependents(ctx context.Context, issueID string) ([]*types.Issue, error) {
	var result []*types.Issue
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetDependentsInTx(ctx, regularTx, issueID)
		return err
	})
	return result, err
}

func (s *Store) GetDependencyTree(ctx context.Context, issueID string, maxDepth int, showAllPaths bool, reverse bool) ([]*types.TreeNode, error) {
	var result []*types.TreeNode
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetDependencyTreeInTx(ctx, regularTx, issueID, maxDepth, showAllPaths, reverse)
		return err
	})
	return result, err
}

func (s *Store) GetIssuesByLabel(ctx context.Context, label string) ([]*types.Issue, error) {
	var ids []string
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		ids, err = issueops.GetIssuesByLabelInTx(ctx, regularTx, label)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetIssuesByIDs(ctx, ids)
}

func (s *Store) GetBlockedIssues(ctx context.Context, filter types.WorkFilter) ([]*types.BlockedIssue, error) {
	var result []*types.BlockedIssue
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetBlockedIssuesInTx(ctx, regularTx, filter)
		return err
	})
	return result, err
}

func (s *Store) GetEpicsEligibleForClosure(ctx context.Context) ([]*types.EpicStatus, error) {
	var result []*types.EpicStatus
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetEpicsEligibleForClosureInTx(ctx, regularTx)
		return err
	})
	return result, err
}

func (s *Store) AddIssueComment(ctx context.Context, issueID, author, text string) (*types.Comment, error) {
	var result *types.Comment
	err := s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.AddIssueCommentInTx(ctx, regularTx, issueID, author, text)
		return err
	})
	return result, err
}

func (s *Store) GetIssueComments(ctx context.Context, issueID string) ([]*types.Comment, error) {
	var result []*types.Comment
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetIssueCommentsInTx(ctx, regularTx, issueID)
		return err
	})
	return result, err
}

func (s *Store) GetEvents(ctx context.Context, issueID string, limit int) ([]*types.Event, error) {
	var result []*types.Event
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetEventsInTx(ctx, regularTx, issueID, limit)
		return err
	})
	return result, err
}

func (s *Store) GetAllEventsSince(ctx context.Context, since time.Time) ([]*types.Event, error) {
	var result []*types.Event
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetAllEventsSinceInTx(ctx, regularTx, since)
		return err
	})
	return result, err
}

func (s *Store) DeleteIssues(ctx context.Context, ids []string, cascade bool, force bool, dryRun bool) (*types.DeleteIssuesResult, error) {
	var result *types.DeleteIssuesResult
	err := s.withConn(ctx, !dryRun, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.DeleteIssuesInTx(ctx, regularTx, ignoredTx, ids, cascade, force, dryRun)
		return err
	})
	return result, err
}

func (s *Store) DeleteIssuesBySourceRepo(ctx context.Context, sourceRepo string) (int, error) {
	var count int
	err := s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		count, err = issueops.DeleteIssuesBySourceRepoInTx(ctx, regularTx, sourceRepo)
		return err
	})
	return count, err
}

func (s *Store) UpdateIssueID(ctx context.Context, oldID, newID string, issue *types.Issue, actor string) error {
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		return issueops.UpdateIssueIDInTx(ctx, regularTx, ignoredTx, oldID, newID, issue, actor)
	})
}

func (s *Store) PromoteFromEphemeral(ctx context.Context, id string, actor string) error {
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		return issueops.PromoteFromEphemeralInTx(ctx, regularTx, ignoredTx, id, actor)
	})
}

func (s *Store) RenameCounterPrefix(ctx context.Context, oldPrefix, newPrefix string) error {
	return nil
}

func (s *Store) RenameDependencyPrefix(ctx context.Context, oldPrefix, newPrefix string) error {
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		return issueops.RenameDependencyPrefixInTx(ctx, regularTx, oldPrefix, newPrefix)
	})
}

func (s *Store) GetDependencyRecords(ctx context.Context, issueID string) ([]*types.Dependency, error) {
	var result []*types.Dependency
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		m, err := issueops.GetDependencyRecordsForIssuesInTx(ctx, regularTx, []string{issueID})
		if err != nil {
			return err
		}
		result = m[issueID]
		return nil
	})
	return result, err
}

func (s *Store) FindWispDependentsRecursive(ctx context.Context, ids []string) (map[string]bool, error) {
	var result map[string]bool
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.FindWispDependentsRecursiveInTx(ctx, ignoredTx, ids)
		return err
	})
	return result, err
}

func (s *Store) AddComment(ctx context.Context, issueID, actor, comment string) error {
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		return issueops.AddCommentEventInTx(ctx, regularTx, issueID, actor, comment)
	})
}

func (s *Store) ImportIssueComment(ctx context.Context, issueID, author, text string, createdAt time.Time) (*types.Comment, error) {
	var result *types.Comment
	err := s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.ImportIssueCommentInTx(ctx, regularTx, issueID, author, text, createdAt)
		return err
	})
	return result, err
}

func (s *Store) GetCommentsForIssues(ctx context.Context, issueIDs []string) (map[string][]*types.Comment, error) {
	var result map[string][]*types.Comment
	err := s.withConn(ctx, false, func(regularTx, ignoredTx *sql.Tx) error {
		var err error
		result, err = issueops.GetCommentsForIssuesInTx(ctx, regularTx, issueIDs)
		return err
	})
	return result, err
}

func (s *Store) ImportJSONLData(
	ctx context.Context,
	issues []*types.Issue,
	configEntries map[string]string,
	actor string,
) (int, error) {
	var imported int
	err := s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		stats := &types.Statistics{}
		if err := issueops.ScanIssueCountsInTx(ctx, regularTx, stats); err != nil {
			return fmt.Errorf("checking issue count: %w", err)
		}
		if stats.TotalIssues > 0 {
			return nil
		}
		for key, value := range configEntries {
			if err := issueops.SetConfigInTx(ctx, regularTx, key, value); err != nil {
				return fmt.Errorf("importing config %q: %w", key, err)
			}
		}
		if len(issues) == 0 {
			return nil
		}
		if _, hasPrefix := configEntries["issue_prefix"]; !hasPrefix {
			firstPrefix := utils.ExtractIssuePrefix(issues[0].ID)
			if firstPrefix != "" {
				if err := issueops.SetConfigInTx(ctx, regularTx, "issue_prefix", firstPrefix); err != nil {
					return fmt.Errorf("setting issue_prefix: %w", err)
				}
			}
		}
		if err := issueops.CreateIssuesInTx(ctx, regularTx, ignoredTx, issues, actor, storage.BatchCreateOptions{
			OrphanHandling:       storage.OrphanAllow,
			SkipPrefixValidation: true,
		}); err != nil {
			return err
		}
		imported = len(issues)
		return nil
	})
	return imported, err
}
