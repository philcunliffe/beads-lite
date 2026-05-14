//go:build sqlite_lite

package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
)

func TestStoreIssueLifecycle(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, filepath.Join(t.TempDir(), "beads.sqlite3"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	if err := store.SetConfig(ctx, "issue_prefix", "lite"); err != nil {
		t.Fatalf("SetConfig(issue_prefix): %v", err)
	}

	blocker := &types.Issue{
		ID:        "lite-blocker",
		Title:     "Blocker",
		Status:    types.StatusOpen,
		Priority:  1,
		IssueType: types.TypeTask,
	}
	task := &types.Issue{
		ID:        "lite-task",
		Title:     "Blocked task",
		Status:    types.StatusOpen,
		Priority:  2,
		IssueType: types.TypeTask,
		Labels:    []string{"sqlite-lite"},
	}
	if err := store.CreateIssue(ctx, blocker, "tester"); err != nil {
		t.Fatalf("CreateIssue(blocker): %v", err)
	}
	if err := store.CreateIssue(ctx, task, "tester"); err != nil {
		t.Fatalf("CreateIssue(task): %v", err)
	}
	if err := store.AddDependency(ctx, &types.Dependency{
		IssueID:     task.ID,
		DependsOnID: blocker.ID,
		Type:        types.DepBlocks,
	}, "tester"); err != nil {
		t.Fatalf("AddDependency: %v", err)
	}

	ready, err := store.GetReadyWork(ctx, types.WorkFilter{})
	if err != nil {
		t.Fatalf("GetReadyWork: %v", err)
	}
	if hasIssue(ready, task.ID) {
		t.Fatalf("%s should be blocked before %s closes", task.ID, blocker.ID)
	}
	if !hasIssue(ready, blocker.ID) {
		t.Fatalf("%s should be ready before it closes", blocker.ID)
	}

	if _, err := store.AddIssueComment(ctx, task.ID, "tester", "works on sqlite"); err != nil {
		t.Fatalf("AddIssueComment: %v", err)
	}
	comments, err := store.GetIssueComments(ctx, task.ID)
	if err != nil {
		t.Fatalf("GetIssueComments: %v", err)
	}
	if len(comments) != 1 || comments[0].Text != "works on sqlite" {
		t.Fatalf("comments = %#v, want one sqlite comment", comments)
	}

	labels, err := store.GetLabelsForIssues(ctx, []string{task.ID})
	if err != nil {
		t.Fatalf("GetLabelsForIssues: %v", err)
	}
	if got := labels[task.ID]; len(got) != 1 || got[0] != "sqlite-lite" {
		t.Fatalf("labels[%s] = %#v, want [sqlite-lite]", task.ID, got)
	}

	if err := store.CloseIssue(ctx, blocker.ID, "done", "tester", ""); err != nil {
		t.Fatalf("CloseIssue(blocker): %v", err)
	}
	ready, err = store.GetReadyWork(ctx, types.WorkFilter{})
	if err != nil {
		t.Fatalf("GetReadyWork after close: %v", err)
	}
	if !hasIssue(ready, task.ID) {
		t.Fatalf("%s should be ready after %s closes", task.ID, blocker.ID)
	}
}

func TestStoreRejectsVersionedCapabilities(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, filepath.Join(t.TempDir(), "beads.sqlite3"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	for name, err := range map[string]error{
		"push":    store.Push(ctx),
		"pull":    store.Pull(ctx),
		"merge":   func() error { _, err := store.Merge(ctx, "main"); return err }(),
		"history": func() error { _, err := store.History(ctx, "lite-task"); return err }(),
		"diff":    func() error { _, err := store.Diff(ctx, "HEAD~1", "HEAD"); return err }(),
	} {
		if !errors.Is(err, storage.ErrUnsupportedCapability) {
			t.Fatalf("%s error = %v, want ErrUnsupportedCapability", name, err)
		}
	}
}

func hasIssue(issues []*types.Issue, id string) bool {
	for _, issue := range issues {
		if issue.ID == id {
			return true
		}
	}
	return false
}
