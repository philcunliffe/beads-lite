//go:build sqlite_lite

package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSQLiteLiteCutsVersionedCommands(t *testing.T) {
	cutCommands := [][]string{
		{"backup"},
		{"branch"},
		{"compact"},
		{"diff"},
		{"dolt"},
		{"federation"},
		{"flatten"},
		{"history"},
		{"merge-slot"},
		{"restore"},
		{"vc"},
	}

	for _, path := range cutCommands {
		if commandPathExists(rootCmd, path) {
			t.Fatalf("sqlite-lite build should not register command path %v", path)
		}
	}
}

func TestSQLiteLiteKeepsLocalMaintenanceCommands(t *testing.T) {
	if !commandPathExists(rootCmd, []string{"gc"}) {
		t.Fatal("sqlite-lite build should register local gc command")
	}
	compact, ok := findCommandPath(rootCmd, []string{"admin", "compact"})
	if !ok {
		t.Fatal("sqlite-lite build should keep admin compact")
	}
	if flag := compact.Flags().Lookup("dolt"); flag != nil {
		t.Fatal("sqlite-lite admin compact should not expose --dolt")
	}
	if commandPathExists(rootCmd, []string{"admin", "compact", "compact"}) {
		t.Fatal("sqlite-lite build should not register root compact-dolt command")
	}
}

func TestSQLiteLiteHidesDoltOnlyGlobalFlags(t *testing.T) {
	for _, name := range []string{"dolt-auto-commit", "sandbox", "global"} {
		if flag := rootCmd.PersistentFlags().Lookup(name); flag != nil {
			t.Fatalf("sqlite-lite build should not expose --%s", name)
		}
	}
	dbFlag := rootCmd.PersistentFlags().Lookup("db")
	if dbFlag == nil {
		t.Fatal("sqlite-lite build should expose --db")
	}
	if dbFlag.Usage != "SQLite database path (default: auto-discover .beads/beads.sqlite3)" {
		t.Fatalf("--db usage = %q", dbFlag.Usage)
	}
}

func TestSQLiteLiteHidesHistoricalShowFlag(t *testing.T) {
	show, ok := findCommandPath(rootCmd, []string{"show"})
	if !ok {
		t.Fatal("sqlite-lite build should register show command")
	}
	if flag := show.Flags().Lookup("as-of"); flag != nil {
		t.Fatal("sqlite-lite show should not expose --as-of")
	}
}

func commandPathExists(root *cobra.Command, path []string) bool {
	_, ok := findCommandPath(root, path)
	return ok
}

func findCommandPath(root *cobra.Command, path []string) (*cobra.Command, bool) {
	cmd := root
	for _, name := range path {
		var next *cobra.Command
		for _, child := range cmd.Commands() {
			if child.Name() == name {
				next = child
				break
			}
			for _, alias := range child.Aliases {
				if alias == name {
					next = child
					break
				}
			}
			if next != nil {
				break
			}
		}
		if next == nil {
			return nil, false
		}
		cmd = next
	}
	return cmd, true
}
