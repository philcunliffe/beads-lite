package fix

import (
	"path/filepath"

	"github.com/steveyegge/beads/internal/configfile"
)

// getDatabasePath returns the database directory recorded in metadata.json,
// falling back to the legacy dolt/ path when no config can be loaded.
func getDatabasePath(beadsDir string) string {
	cfg, err := configfile.Load(beadsDir)
	if err != nil || cfg == nil {
		return filepath.Join(beadsDir, "dolt")
	}
	return cfg.DatabasePath(beadsDir)
}
