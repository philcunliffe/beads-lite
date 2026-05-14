//go:build sqlite_lite

package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/internal/types"
)

var (
	gcDryRun    bool
	gcForce     bool
	gcOlderThan int
	gcSkipDecay bool
)

var gcCmd = &cobra.Command{
	Use:     "gc",
	GroupID: "maint",
	Short:   "Delete old closed issues from the local SQLite database",
	Long: `Garbage collect closed local issues.

SQLite lite has no Dolt history, remotes, or engine GC phase. This command only
removes closed, unpinned issues older than the selected retention window.`,
	Run: func(cmd *cobra.Command, _ []string) {
		if !gcDryRun {
			CheckReadonly("gc")
		}
		if gcOlderThan < 0 {
			FatalError("--older-than must be non-negative")
		}
		start := time.Now()
		deleted := 0
		candidates := 0

		if !gcSkipDecay {
			cutoff := time.Now().AddDate(0, 0, -gcOlderThan)
			statusClosed := types.StatusClosed
			issues, err := store.SearchIssues(rootCtx, "", types.IssueFilter{
				Status:       &statusClosed,
				ClosedBefore: &cutoff,
			})
			if err != nil {
				FatalError("searching closed issues: %v", err)
			}
			for _, issue := range issues {
				if issue.Pinned {
					continue
				}
				candidates++
				if gcDryRun {
					continue
				}
				if !gcForce {
					FatalErrorWithHint(
						fmt.Sprintf("would delete %d closed issue(s) older than %d days", candidates, gcOlderThan),
						"Use --force to confirm or --dry-run to preview.")
				}
				if err := store.DeleteIssue(rootCtx, issue.ID); err != nil {
					WarnError("failed to delete %s: %v", issue.ID, err)
					continue
				}
				deleted++
			}
			if deleted > 0 {
				commandDidWrite.Store(true)
			}
		}

		if jsonOutput {
			outputJSON(map[string]interface{}{
				"dry_run":    gcDryRun,
				"candidates": candidates,
				"deleted":    deleted,
				"elapsed_ms": time.Since(start).Milliseconds(),
			})
			return
		}
		if gcSkipDecay {
			fmt.Println("SQLite GC: decay skipped")
			return
		}
		if gcDryRun {
			fmt.Printf("SQLite GC dry run: would delete %d issue(s)\n", candidates)
			return
		}
		fmt.Printf("SQLite GC complete: deleted %d issue(s)\n", deleted)
	},
}

func init() {
	gcCmd.Flags().BoolVar(&gcDryRun, "dry-run", false, "Preview without making changes")
	gcCmd.Flags().BoolVarP(&gcForce, "force", "f", false, "Skip confirmation prompts")
	gcCmd.Flags().IntVar(&gcOlderThan, "older-than", 90, "Delete closed issues older than N days")
	gcCmd.Flags().BoolVar(&gcSkipDecay, "skip-decay", false, "Skip issue deletion phase")
	rootCmd.AddCommand(gcCmd)
}
