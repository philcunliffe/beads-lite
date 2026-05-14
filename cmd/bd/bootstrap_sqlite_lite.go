//go:build sqlite_lite

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/internal/beads"
	"github.com/steveyegge/beads/internal/config"
	"github.com/steveyegge/beads/internal/configfile"
)

var bootstrapCmd = &cobra.Command{
	Use:     "bootstrap",
	GroupID: "setup",
	Short:   "Create or import a local SQLite beads database",
	RunE: func(cmd *cobra.Command, args []string) error {
		beadsDir := beads.FindBeadsDir()
		if beadsDir == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			beadsDir = filepath.Join(cwd, ".beads")
		}
		if _, err := os.Stat(filepath.Join(beadsDir, "metadata.json")); os.IsNotExist(err) {
			if err := os.MkdirAll(beadsDir, config.BeadsDirPerm); err != nil {
				return fmt.Errorf("creating .beads directory: %w", err)
			}
			cwd, _ := os.Getwd()
			prefix := sanitizeLitePrefix(filepath.Base(cwd))
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
			if err := st.SetConfig(rootCtx, "issue_prefix", prefix); err != nil {
				_ = st.Close()
				return fmt.Errorf("setting issue prefix: %w", err)
			}
			_ = st.Close()
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
		jsonlPath := filepath.Join(beadsDir, "issues.jsonl")
		if info, err := os.Stat(jsonlPath); err == nil && info.Size() > 0 {
			f, err := os.Open(jsonlPath)
			if err != nil {
				return err
			}
			defer f.Close()
			return runImportFromReader(rootCtx, f, jsonlPath)
		}
		if !quietFlag {
			fmt.Fprintf(os.Stderr, "SQLite beads database ready in %s\n", beadsDir)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
}
