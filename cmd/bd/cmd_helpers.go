package main

import (
	"context"
	"os/exec"

	"github.com/steveyegge/beads/internal/storage"
)

// transact runs fn inside the store's transaction and commits with msg.
// Restored helper after the legacy Dolt-mode helpers were stripped.
func transact(ctx context.Context, s storage.DoltStorage, msg string, fn func(tx storage.Transaction) error) error {
	return s.RunInTransaction(ctx, msg, fn)
}

// isDoltNothingToCommit historically detected the "nothing to commit" error
// returned by the Dolt engine. In lite mode SQLite never produces it, so we
// keep the predicate for source compatibility and always return false.
func isDoltNothingToCommit(_ error) bool { return false }

// doltAutoCommitParams collected the configuration the now-removed auto-commit
// path needed. The struct is preserved so existing call sites compile; the
// values are ignored in lite mode.
type doltAutoCommitParams struct {
	Command         string
	IssueIDs        []string
	Actor           string
	Issue           string
	Message         string
	MessageOverride string
	Reason          string
}

const (
	doltAutoCommitOn    = "on"
	doltAutoCommitOff   = "off"
	doltAutoCommitBatch = "batch"
)

// maybeAutoCommit is the legacy auto-commit dispatcher. The lite build commits
// inline through SQLite, so this is a no-op.
func maybeAutoCommit(_ context.Context, _ doltAutoCommitParams) error { return nil }

// formatDoltAutoCommitMessage formatted the auto-commit message under the
// legacy Dolt path. Kept for signature compatibility.
func formatDoltAutoCommitMessage(cmdName, actor string, _ []string) string {
	if cmdName == "" {
		return actor
	}
	return cmdName + ": " + actor
}

// maybeAutoExport ran the JSONL auto-export under the legacy Dolt path. No-op
// in lite mode (writes are already in the local SQLite file).
func maybeAutoExport(_ context.Context) {}

// maybeAutoPush pushed pending Dolt commits to the configured remote. No-op
// in lite mode.
func maybeAutoPush(_ context.Context) {}

// versionChange mirrors the entry shape the legacy upgrade printer expected.
type versionChange struct {
	Version string
	Date    string
	Changes []string
}

// getVersionsSince walked the changelog to compute which versions had been
// applied since the last upgrade hint. Lite mode has no curated changelog,
// so always returns an empty list.
func getVersionsSince(_ string) []versionChange { return nil }

// maybeAutoCommitStore was the auto-commit dispatch for Dolt-backed stores.
// In lite mode SQLite commits inline, so this is a no-op.
func maybeAutoCommitStore(_ context.Context, _ storage.DoltStorage, _ doltAutoCommitParams) error {
	return nil
}

// configSideEffect describes a side-effect that a config change would have
// had under the legacy Dolt mode. The lite build retains the type and
// helpers so config command code compiles, but never reports any effects.
type configSideEffect struct {
	Description string
	Recommended string
}

func checkConfigSetSideEffects(_ string, _ string) []configSideEffect { return nil }
func checkConfigUnsetSideEffects(_ string) []configSideEffect         { return nil }
func printConfigSideEffects(_ []configSideEffect)                     {}

// applyFixList was the legacy dispatch for doctor --fix follow-ups. The lite
// build has no fixable doctor checks, so applying a fix list is a no-op.
func applyFixList(_ string, _ []doctorCheck) {}

// gitAddFile re-adds a file to the git index. The lite build retains the
// helper because the export-auto hook code paths still want to opportunistically
// stage exported JSONL files. It is best-effort: failures from a non-git
// directory are returned as-is so callers can decide whether to ignore them.
func gitAddFile(path string) error {
	return runQuietGit("add", path)
}

func runQuietGit(args ...string) error {
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

// maybeShowUpgradeNotification used to print a hint when a newer bd was
// available. The lite build skips remote version checks entirely.
func maybeShowUpgradeNotification() {}

// getDoltAutoCommitMode returns the legacy Dolt auto-commit mode. In lite
// mode commits happen inline; returning "off" keeps callers happy.
func getDoltAutoCommitMode() (string, error) { return "off", nil }

// maybeWarnLinearStaleness was the Linear-sync staleness banner. Lite has no
// Linear integration.
func maybeWarnLinearStaleness(_ interface{}) {}

// isSandboxed reports whether bd is running inside a sandboxed environment
// that should block writes. The lite build conservatively returns false.
func isSandboxed() bool { return false }
