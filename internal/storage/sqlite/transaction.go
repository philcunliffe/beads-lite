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
)

// RunInTransaction executes a function within a single SQLite transaction.
func (s *Store) RunInTransaction(ctx context.Context, commitMsg string, fn func(tx storage.Transaction) error) error {
	_ = commitMsg
	return s.withConn(ctx, true, func(regularTx, ignoredTx *sql.Tx) error {
		tx := &sqliteTransaction{regularTx: regularTx, ignoredTx: ignoredTx}
		return fn(tx)
	})
}

type sqliteTransaction struct {
	regularTx *sql.Tx
	ignoredTx *sql.Tx
}

func (t *sqliteTransaction) CreateIssue(ctx context.Context, issue *types.Issue, actor string) error {
	bc, err := issueops.NewBatchContext(ctx, t.regularTx, storage.BatchCreateOptions{SkipPrefixValidation: true})
	if err != nil {
		return err
	}
	return issueops.CreateIssueInTx(ctx, t.regularTx, t.ignoredTx, bc, issue, actor)
}

func (t *sqliteTransaction) CreateIssues(ctx context.Context, issues []*types.Issue, actor string) error {
	for _, issue := range issues {
		if err := t.CreateIssue(ctx, issue, actor); err != nil {
			return err
		}
	}
	return nil
}

func (t *sqliteTransaction) UpdateIssue(ctx context.Context, id string, updates map[string]interface{}, actor string) error {
	_, err := issueops.UpdateIssueInTx(ctx, t.regularTx, id, updates, actor)
	return err
}

func (t *sqliteTransaction) CloseIssue(ctx context.Context, id string, reason string, actor string, session string) error {
	_, err := issueops.CloseIssueInTx(ctx, t.regularTx, id, reason, actor, session)
	return err
}

func (t *sqliteTransaction) DeleteIssue(ctx context.Context, id string) error {
	return issueops.DeleteIssueInTx(ctx, t.regularTx, t.ignoredTx, id)
}

func (t *sqliteTransaction) GetIssue(ctx context.Context, id string) (*types.Issue, error) {
	return issueops.GetIssueInTx(ctx, t.regularTx, id)
}

func (t *sqliteTransaction) SearchIssues(ctx context.Context, query string, filter types.IssueFilter) ([]*types.Issue, error) {
	return issueops.SearchIssuesInTx(ctx, t.regularTx, query, filter)
}

func (t *sqliteTransaction) AddDependency(ctx context.Context, dep *types.Dependency, actor string) error {
	return t.AddDependencyWithOptions(ctx, dep, actor, storage.DependencyAddOptions{})
}

func (t *sqliteTransaction) AddDependencyWithOptions(ctx context.Context, dep *types.Dependency, actor string, addOpts storage.DependencyAddOptions) error {
	return issueops.AddDependencyInTx(ctx, t.regularTx, dep, actor, issueops.AddDependencyOpts{
		IsCrossPrefix:  types.ExtractPrefix(dep.IssueID) != types.ExtractPrefix(dep.DependsOnID),
		SkipCycleCheck: addOpts.SkipCycleCheck,
	})
}

func (t *sqliteTransaction) RemoveDependency(ctx context.Context, issueID, dependsOnID string, actor string) error {
	return issueops.RemoveDependencyInTx(ctx, t.regularTx, issueID, dependsOnID)
}

func (t *sqliteTransaction) GetDependencyRecords(ctx context.Context, issueID string) ([]*types.Dependency, error) {
	m, err := issueops.GetDependencyRecordsForIssuesInTx(ctx, t.regularTx, []string{issueID})
	if err != nil {
		return nil, err
	}
	return m[issueID], nil
}

func (t *sqliteTransaction) AddLabel(ctx context.Context, issueID, label, actor string) error {
	return issueops.AddLabelInTx(ctx, t.regularTx, "", "", issueID, label, actor)
}

func (t *sqliteTransaction) RemoveLabel(ctx context.Context, issueID, label, actor string) error {
	return issueops.RemoveLabelInTx(ctx, t.regularTx, "", "", issueID, label, actor)
}

func (t *sqliteTransaction) GetLabels(ctx context.Context, issueID string) ([]string, error) {
	return issueops.GetLabelsInTx(ctx, t.regularTx, "", issueID)
}

func (t *sqliteTransaction) SetConfig(ctx context.Context, key, value string) error {
	if err := issueops.SetConfigInTx(ctx, t.regularTx, key, value); err != nil {
		return err
	}
	// Sync normalized tables when config keys change
	switch key {
	case "status.custom":
		if err := issueops.SyncCustomStatusesTable(ctx, t.regularTx, value); err != nil {
			return fmt.Errorf("syncing custom_statuses table: %w", err)
		}
	case "types.custom":
		if err := issueops.SyncCustomTypesTable(ctx, t.regularTx, value); err != nil {
			return fmt.Errorf("syncing custom_types table: %w", err)
		}
	}
	return nil
}

func (t *sqliteTransaction) GetConfig(ctx context.Context, key string) (string, error) {
	return issueops.GetConfigInTx(ctx, t.regularTx, key)
}

func (t *sqliteTransaction) SetMetadata(ctx context.Context, key, value string) error {
	return issueops.SetMetadataInTx(ctx, t.regularTx, key, value)
}

func (t *sqliteTransaction) GetMetadata(ctx context.Context, key string) (string, error) {
	return issueops.GetMetadataInTx(ctx, t.regularTx, key)
}

func (t *sqliteTransaction) SetLocalMetadata(ctx context.Context, key, value string) error {
	return issueops.SetLocalMetadataInTx(ctx, t.ignoredTx, key, value)
}

func (t *sqliteTransaction) GetLocalMetadata(ctx context.Context, key string) (string, error) {
	return issueops.GetLocalMetadataInTx(ctx, t.ignoredTx, key)
}

func (t *sqliteTransaction) AddComment(ctx context.Context, issueID, actor, comment string) error {
	return fmt.Errorf("sqliteTransaction: AddComment not implemented")
}

func (t *sqliteTransaction) ImportIssueComment(ctx context.Context, issueID, author, text string, createdAt time.Time) (*types.Comment, error) {
	return nil, fmt.Errorf("sqliteTransaction: ImportIssueComment not implemented")
}

func (t *sqliteTransaction) GetIssueComments(ctx context.Context, issueID string) ([]*types.Comment, error) {
	return nil, fmt.Errorf("sqliteTransaction: GetIssueComments not implemented")
}

func (t *sqliteTransaction) CreateIssueImport(ctx context.Context, issue *types.Issue, actor string, skipPrefixValidation bool) error {
	bc, err := issueops.NewBatchContext(ctx, t.regularTx, storage.BatchCreateOptions{SkipPrefixValidation: skipPrefixValidation})
	if err != nil {
		return err
	}
	return issueops.CreateIssueInTx(ctx, t.regularTx, t.ignoredTx, bc, issue, actor)
}
