package doctor

// RunPerformanceDiagnostics is a lite-mode no-op (the SQLite build has no
// performance diagnostics yet). Returning nil keeps the --perf flag usable
// without surprising callers.
func RunPerformanceDiagnostics(_ string) error { return nil }

// CollectPlatformInfo returns a minimal platform map for the export path.
func CollectPlatformInfo(_ string) map[string]string {
	return map[string]string{}
}

// CheckFreshClone is a lite-mode no-op check; fresh-clone detection lived in
// the (now-removed) Dolt server path.
func CheckFreshClone(_ string) DoctorCheck {
	return liteSkip("Fresh Clone")
}

// CheckDatabaseConfig was a Dolt-server config validation step.
func CheckDatabaseConfig(_ string) DoctorCheck {
	return liteSkip("Database Config")
}

// CheckMultiRepoTypes validated multi-repo type metadata.
func CheckMultiRepoTypes(_ string) DoctorCheck {
	return liteSkip("Multi-Repo Types")
}

// CheckBeadsRoleWithStore validated the maintainer/contributor role.
func CheckBeadsRoleWithStore(_ string, _ *SharedStore) DoctorCheck {
	return liteSkip("Beads Role")
}

// CheckStaleLockFiles checked for stale Dolt server lock files.
func CheckStaleLockFiles(_ string) DoctorCheck {
	return liteSkip("Stale Lock Files")
}

// CheckRemoteConsistency validated Dolt remote consistency.
func CheckRemoteConsistency(_ string) DoctorCheck {
	return liteSkip("Remote Consistency")
}

// RunDoltHealthChecks ran the full Dolt server health battery.
func RunDoltHealthChecks(_ string) []DoctorCheck {
	return nil
}

// CheckDoltConnection verified a live Dolt connection.
func CheckDoltConnection(_ string) DoctorCheck {
	return liteSkip("Dolt Connection")
}

// CheckDoltSchema verified the Dolt schema version.
func CheckDoltSchema(_ string) DoctorCheck {
	return liteSkip("Dolt Schema")
}

// CheckDoltLocks looked at the Dolt server lock state.
func CheckDoltLocks(_ string) DoctorCheck {
	return liteSkip("Dolt Locks")
}

// GetSuppressedChecksWithStore returned the set of checks the user opted out
// of in config; in lite mode no checks are suppressed.
func GetSuppressedChecksWithStore(_ string, _ *SharedStore) map[string]bool {
	return map[string]bool{}
}

// CheckKVSyncStatus inspected the Dolt KV sync state.
func CheckKVSyncStatus(_ string) DoctorCheck {
	return liteSkip("KV Sync")
}

// CheckBtrfsNoCOW verified the btrfs NoCOW attribute on .beads/.
func CheckBtrfsNoCOW(_ string) DoctorCheck {
	return liteSkip("Btrfs NoCOW")
}

// CheckChildParentDependencies scanned for child→parent dependency anti-patterns.
func CheckChildParentDependencies(_ string) DoctorCheck {
	return liteSkip("Child-Parent Dependencies")
}

// CheckDuplicateIssues detected duplicate issues across multi-repo hydration.
func CheckDuplicateIssues(_ string, _ bool, _ int) DoctorCheck {
	return liteSkip("Duplicate Issues")
}

// CheckStaleMolecules detected stale molecule rows.
func CheckStaleMolecules(_ string) DoctorCheck {
	return liteSkip("Stale Molecules")
}

// CheckPersistentMolIssues looked for orphaned mol- prefix issues.
func CheckPersistentMolIssues(_ string) DoctorCheck {
	return liteSkip("Persistent Mol Issues")
}

// CheckStaleMQFiles detected stale message queue files.
func CheckStaleMQFiles(_ string) DoctorCheck {
	return liteSkip("Stale MQ Files")
}

// CheckStaleClosedIssues detected stale closed issues left in molecules.
func CheckStaleClosedIssues(_ string) DoctorCheck {
	return liteSkip("Stale Closed Issues")
}

// CheckPatrolPollution detected patrol pollution in molecules.
func CheckPatrolPollution(_ string) DoctorCheck {
	return liteSkip("Patrol Pollution")
}

// CheckTestPollution detected test pollution in molecules.
func CheckTestPollution(_ string) DoctorCheck {
	return liteSkip("Test Pollution")
}

// CheckNameToSlug converts a doctor check name into a config-key slug. The
// lite implementation keeps the historical convention: lowercase, spaces
// turned into dashes, other punctuation stripped.
func CheckNameToSlug(name string) string {
	var b []byte
	for _, r := range name {
		switch {
		case r >= 'A' && r <= 'Z':
			b = append(b, byte(r-'A'+'a'))
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b = append(b, byte(r))
		case r == ' ', r == '-', r == '_':
			if len(b) > 0 && b[len(b)-1] != '-' {
				b = append(b, '-')
			}
		}
	}
	for len(b) > 0 && b[len(b)-1] == '-' {
		b = b[:len(b)-1]
	}
	return string(b)
}

// CheckMigrationReadiness inspected pending migrations.
func CheckMigrationReadiness(_ string) (DoctorCheck, MigrationValidationResult) {
	return liteSkip("Migration Readiness"), MigrationValidationResult{Ready: true, Backend: "sqlite"}
}

// CheckMigrationCompletion verified completed migrations.
func CheckMigrationCompletion(_ string) (DoctorCheck, MigrationValidationResult) {
	return liteSkip("Migration Completion"), MigrationValidationResult{Ready: true, Backend: "sqlite"}
}

// MigrationValidationResult was the legacy migration-validation report.
type MigrationValidationResult struct {
	OK             bool
	Ready          bool
	Backend        string
	JSONLCount     int
	SQLiteCount    int
	DoltCount      int
	JSONLValid     bool
	JSONLMalformed int
	Errors         []string
	Warnings       []string
}

// CheckFederationRemotesAPI ran an API-level federation remotes audit.
func CheckFederationRemotesAPI(_ string) DoctorCheck {
	return liteSkip("Federation Remotes API")
}

// CheckFederationPeerConnectivity reached out to each federation peer.
func CheckFederationPeerConnectivity(_ string) DoctorCheck {
	return liteSkip("Federation Peer Connectivity")
}

// CheckFederationSyncStaleness measured federation sync staleness.
func CheckFederationSyncStaleness(_ string) DoctorCheck {
	return liteSkip("Federation Sync Staleness")
}

// CheckFederationConflicts surfaced unresolved federation conflicts.
func CheckFederationConflicts(_ string) DoctorCheck {
	return liteSkip("Federation Conflicts")
}

// CheckDoltServerModeMismatch flagged Dolt server-mode misconfiguration.
func CheckDoltServerModeMismatch(_ string) DoctorCheck {
	return liteSkip("Dolt Server Mode Mismatch")
}

// CheckAgentDocumentation audited AGENTS.md/CLAUDE.md content.
func CheckAgentDocumentation(_ string) DoctorCheck {
	return liteSkip("Agent Documentation")
}

// CheckAgentDocDivergence checked for divergence across agent doc files.
func CheckAgentDocDivergence(_ string) DoctorCheck {
	return liteSkip("Agent Documentation Divergence")
}

// CheckLegacyBeadsSlashCommands looked for stale beads slash-command files.
func CheckLegacyBeadsSlashCommands(_ string) DoctorCheck {
	return liteSkip("Legacy Beads Slash Commands")
}

// CheckLegacyMCPToolReferences looked for stale MCP tool references.
func CheckLegacyMCPToolReferences(_ string) DoctorCheck {
	return liteSkip("Legacy MCP Tool References")
}

// CheckOrphanedDependencies scanned for orphaned dependency rows.
func CheckOrphanedDependencies(_ string) DoctorCheck {
	return liteSkip("Orphaned Dependencies")
}

// CheckGitConflicts looked for unresolved git conflict markers.
func CheckGitConflicts(_ string) DoctorCheck {
	return liteSkip("Git Conflicts")
}

func liteSkip(name string) DoctorCheck {
	return DoctorCheck{
		Name:    name,
		Status:  StatusOK,
		Message: "Not applicable in lite mode",
	}
}
