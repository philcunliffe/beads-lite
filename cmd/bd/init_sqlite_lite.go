//go:build sqlite_lite

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/internal/config"
	"github.com/steveyegge/beads/internal/configfile"
)

var initCmd = &cobra.Command{
	Use:     "init",
	GroupID: "setup",
	Short:   "Initialize a local SQLite beads database",
	RunE: func(cmd *cobra.Command, args []string) error {
		prefix, _ := cmd.Flags().GetString("prefix")
		quiet, _ := cmd.Flags().GetBool("quiet")
		force, _ := cmd.Flags().GetBool("force")
		fromJSONL, _ := cmd.Flags().GetString("from-jsonl")

		beadsDir, err := sqliteLiteInitBeadsDir()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(beadsDir, config.BeadsDirPerm); err != nil {
			return fmt.Errorf("creating .beads directory: %w", err)
		}
		dbFile := filepath.Join(beadsDir, sqliteLiteDatabaseFile)
		if _, err := os.Stat(dbFile); err == nil && !force {
			return fmt.Errorf("database already exists at %s (use --force to reuse it)", dbFile)
		}

		if strings.TrimSpace(prefix) == "" {
			cwd, _ := os.Getwd()
			prefix = sanitizeLitePrefix(filepath.Base(cwd))
		}
		cfg := configfile.DefaultConfig()
		cfg.Backend = configfile.BackendSQLite
		cfg.Database = sqliteLiteDatabaseFile
		cfg.ProjectID = configfile.GenerateProjectID()
		if err := cfg.Save(beadsDir); err != nil {
			return fmt.Errorf("writing metadata.json: %w", err)
		}
		if err := createConfigYaml(beadsDir, false, prefix); err != nil {
			return err
		}
		if err := createReadme(beadsDir); err != nil {
			return err
		}

		st, err := openSQLiteLiteStore(rootCtx, beadsDir)
		if err != nil {
			return err
		}
		store = st
		defer func() {
			_ = st.Close()
			store = nil
		}()
		if err := st.SetConfig(rootCtx, "issue_prefix", prefix); err != nil {
			return fmt.Errorf("setting issue prefix: %w", err)
		}
		if fromJSONL != "" {
			f, err := os.Open(fromJSONL) //nolint:gosec // CLI-selected import file
			if err != nil {
				return fmt.Errorf("opening JSONL import: %w", err)
			}
			defer f.Close()
			if err := runImportFromReader(rootCtx, f, fromJSONL); err != nil {
				return err
			}
		}
		if !quiet {
			fmt.Fprintf(os.Stderr, "Initialized SQLite beads database in %s\n", beadsDir)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().String("prefix", "", "Issue ID prefix (default: current directory name)")
	initCmd.Flags().Bool("quiet", false, "Suppress non-essential output")
	initCmd.Flags().Bool("force", false, "Reuse an existing SQLite database")
	initCmd.Flags().String("from-jsonl", "", "Import issues from a JSONL export after initialization")
	// Accept (and ignore) the flags gascity's lite-mode init expects.
	initCmd.Flags().Bool("skip-hooks", false, "Skip installing git hooks (accepted for gascity compatibility; lite mode never installs hooks)")
	initCmd.Flags().Bool("skip-agents", false, "Skip writing AGENTS.md/CLAUDE.md (accepted for gascity compatibility; lite mode never writes them)")
	rootCmd.AddCommand(initCmd)
}

func sqliteLiteInitBeadsDir() (string, error) {
	if dbPath != "" {
		if beadsDir := resolveCommandBeadsDir(dbPath); beadsDir != "" {
			return beadsDir, nil
		}
		absPath, err := filepath.Abs(dbPath)
		if err != nil {
			return "", err
		}
		return filepath.Dir(absPath), nil
	}
	if env := os.Getenv("BEADS_DIR"); env != "" {
		return env, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".beads"), nil
}
