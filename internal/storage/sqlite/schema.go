//go:build sqlite_lite

package sqlite

import (
	"context"
	"fmt"
	"strings"
)

const sqliteSchemaVersion = 1

func (s *Store) initSchema(ctx context.Context) error {
	for _, stmt := range splitSQL(sqliteSchemaSQL) {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("sqlite: schema statement failed: %w\n%s", err, stmt)
		}
	}
	if _, err := s.db.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", sqliteSchemaVersion)); err != nil {
		return fmt.Errorf("sqlite: set schema version: %w", err)
	}
	return nil
}

func splitSQL(sql string) []string {
	parts := strings.Split(sql, ";")
	stmts := make([]string, 0, len(parts))
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt != "" {
			stmts = append(stmts, stmt)
		}
	}
	return stmts
}

const issueColumnsSQL = `
	id TEXT PRIMARY KEY,
	content_hash TEXT,
	title TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	design TEXT NOT NULL DEFAULT '',
	acceptance_criteria TEXT NOT NULL DEFAULT '',
	notes TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'open',
	priority INTEGER NOT NULL DEFAULT 2,
	issue_type TEXT NOT NULL DEFAULT 'task',
	assignee TEXT,
	estimated_minutes INTEGER,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by TEXT DEFAULT '',
	owner TEXT DEFAULT '',
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	started_at DATETIME,
	closed_at DATETIME,
	closed_by_session TEXT DEFAULT '',
	external_ref TEXT,
	spec_id TEXT,
	compaction_level INTEGER DEFAULT 0,
	compacted_at DATETIME,
	compacted_at_commit TEXT,
	original_size INTEGER,
	source_repo TEXT DEFAULT '',
	close_reason TEXT DEFAULT '',
	sender TEXT DEFAULT '',
	ephemeral INTEGER DEFAULT 0,
	no_history INTEGER DEFAULT 0,
	wisp_type TEXT DEFAULT '',
	pinned INTEGER DEFAULT 0,
	is_template INTEGER DEFAULT 0,
	await_type TEXT DEFAULT '',
	await_id TEXT DEFAULT '',
	timeout_ns INTEGER DEFAULT 0,
	waiters TEXT DEFAULT '',
	mol_type TEXT DEFAULT '',
	event_kind TEXT DEFAULT '',
	actor TEXT DEFAULT '',
	target TEXT DEFAULT '',
	payload TEXT DEFAULT '',
	due_at DATETIME,
	defer_until DATETIME,
	work_type TEXT DEFAULT 'mutex',
	source_system TEXT DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}'
`

var sqliteSchemaSQL = `
CREATE TABLE IF NOT EXISTS issues (` + issueColumnsSQL + `);
CREATE TABLE IF NOT EXISTS wisps (` + issueColumnsSQL + `);

CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
CREATE INDEX IF NOT EXISTS idx_issues_priority ON issues(priority);
CREATE INDEX IF NOT EXISTS idx_issues_issue_type ON issues(issue_type);
CREATE INDEX IF NOT EXISTS idx_issues_assignee ON issues(assignee);
CREATE INDEX IF NOT EXISTS idx_issues_created_at ON issues(created_at);
CREATE INDEX IF NOT EXISTS idx_issues_external_ref ON issues(external_ref);
CREATE INDEX IF NOT EXISTS idx_issues_spec_id ON issues(spec_id);
CREATE INDEX IF NOT EXISTS idx_wisps_status ON wisps(status);
CREATE INDEX IF NOT EXISTS idx_wisps_priority ON wisps(priority);
CREATE INDEX IF NOT EXISTS idx_wisps_issue_type ON wisps(issue_type);
CREATE INDEX IF NOT EXISTS idx_wisps_assignee ON wisps(assignee);
CREATE INDEX IF NOT EXISTS idx_wisps_created_at ON wisps(created_at);
CREATE INDEX IF NOT EXISTS idx_wisps_external_ref ON wisps(external_ref);
CREATE INDEX IF NOT EXISTS idx_wisps_spec_id ON wisps(spec_id);

CREATE TABLE IF NOT EXISTS dependencies (
	issue_id TEXT NOT NULL,
	depends_on_id TEXT NOT NULL,
	type TEXT NOT NULL DEFAULT 'blocks',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by TEXT NOT NULL DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}',
	thread_id TEXT DEFAULT '',
	PRIMARY KEY (issue_id, depends_on_id),
	FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_dependencies_issue ON dependencies(issue_id);
CREATE INDEX IF NOT EXISTS idx_dependencies_depends_on ON dependencies(depends_on_id);
CREATE INDEX IF NOT EXISTS idx_dependencies_depends_on_type ON dependencies(depends_on_id, type);
CREATE INDEX IF NOT EXISTS idx_dependencies_thread ON dependencies(thread_id);

CREATE TABLE IF NOT EXISTS labels (
	issue_id TEXT NOT NULL,
	label TEXT NOT NULL,
	PRIMARY KEY (issue_id, label),
	FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_labels_label ON labels(label);

CREATE TABLE IF NOT EXISTS comments (
	id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
	issue_id TEXT NOT NULL,
	author TEXT NOT NULL,
	text TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_comments_issue ON comments(issue_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at);

CREATE TABLE IF NOT EXISTS events (
	id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
	issue_id TEXT NOT NULL,
	event_type TEXT NOT NULL,
	actor TEXT NOT NULL DEFAULT '',
	old_value TEXT DEFAULT '',
	new_value TEXT DEFAULT '',
	comment TEXT DEFAULT '',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_events_issue ON events(issue_id);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);

CREATE TABLE IF NOT EXISTS wisp_dependencies (
	issue_id TEXT NOT NULL,
	depends_on_id TEXT NOT NULL,
	type TEXT NOT NULL DEFAULT 'blocks',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by TEXT NOT NULL DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}',
	thread_id TEXT DEFAULT '',
	PRIMARY KEY (issue_id, depends_on_id)
);
CREATE INDEX IF NOT EXISTS idx_wisp_dep_depends ON wisp_dependencies(depends_on_id);
CREATE INDEX IF NOT EXISTS idx_wisp_dep_depends_type ON wisp_dependencies(depends_on_id, type);

CREATE TABLE IF NOT EXISTS wisp_labels (
	issue_id TEXT NOT NULL,
	label TEXT NOT NULL,
	PRIMARY KEY (issue_id, label)
);
CREATE INDEX IF NOT EXISTS idx_wisp_labels_label ON wisp_labels(label);

CREATE TABLE IF NOT EXISTS wisp_comments (
	id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
	issue_id TEXT NOT NULL,
	author TEXT NOT NULL DEFAULT '',
	text TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_wisp_comments_issue ON wisp_comments(issue_id);

CREATE TABLE IF NOT EXISTS wisp_events (
	id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
	issue_id TEXT NOT NULL,
	event_type TEXT NOT NULL,
	actor TEXT NOT NULL DEFAULT '',
	old_value TEXT DEFAULT '',
	new_value TEXT DEFAULT '',
	comment TEXT DEFAULT '',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_wisp_events_issue ON wisp_events(issue_id);
CREATE INDEX IF NOT EXISTS idx_wisp_events_created_at ON wisp_events(created_at);

CREATE TABLE IF NOT EXISTS config (` + "`key`" + ` TEXT PRIMARY KEY, value TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS metadata (` + "`key`" + ` TEXT PRIMARY KEY, value TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS local_metadata (` + "`key`" + ` TEXT PRIMARY KEY, value TEXT NOT NULL DEFAULT '');
CREATE TABLE IF NOT EXISTS child_counters (
	parent_id TEXT PRIMARY KEY,
	last_child INTEGER NOT NULL DEFAULT 0,
	FOREIGN KEY (parent_id) REFERENCES issues(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS issue_counter (
	prefix TEXT PRIMARY KEY,
	last_id INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS repo_mtimes (
	repo_path TEXT PRIMARY KEY,
	jsonl_path TEXT NOT NULL,
	mtime_ns INTEGER NOT NULL,
	last_checked DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_repo_mtimes_checked ON repo_mtimes(last_checked);

CREATE TABLE IF NOT EXISTS issue_snapshots (
	id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
	issue_id TEXT NOT NULL,
	snapshot_time DATETIME NOT NULL,
	compaction_level INTEGER NOT NULL,
	original_size INTEGER NOT NULL,
	compressed_size INTEGER NOT NULL,
	original_content TEXT NOT NULL,
	archived_events TEXT,
	FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_snapshots_issue ON issue_snapshots(issue_id);
CREATE INDEX IF NOT EXISTS idx_snapshots_level ON issue_snapshots(compaction_level);

CREATE TABLE IF NOT EXISTS compaction_snapshots (
	id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
	issue_id TEXT NOT NULL,
	compaction_level INTEGER NOT NULL,
	snapshot_json BLOB NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_comp_snap_issue ON compaction_snapshots(issue_id, compaction_level, created_at DESC);

CREATE TABLE IF NOT EXISTS routes (
	prefix TEXT PRIMARY KEY,
	path TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS interactions (
	id TEXT PRIMARY KEY,
	kind TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	actor TEXT,
	issue_id TEXT,
	model TEXT,
	prompt TEXT,
	response TEXT,
	error TEXT,
	tool_name TEXT,
	exit_code INTEGER,
	parent_id TEXT,
	label TEXT,
	reason TEXT,
	extra TEXT
);
CREATE INDEX IF NOT EXISTS idx_interactions_kind ON interactions(kind);
CREATE INDEX IF NOT EXISTS idx_interactions_created_at ON interactions(created_at);
CREATE INDEX IF NOT EXISTS idx_interactions_issue_id ON interactions(issue_id);
CREATE INDEX IF NOT EXISTS idx_interactions_parent_id ON interactions(parent_id);

CREATE TABLE IF NOT EXISTS federation_peers (
	name TEXT PRIMARY KEY,
	remote_url TEXT NOT NULL,
	username TEXT,
	password_encrypted BLOB,
	sovereignty TEXT DEFAULT '',
	last_sync DATETIME,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_federation_peers_sovereignty ON federation_peers(sovereignty);

CREATE TABLE IF NOT EXISTS custom_statuses (
	name TEXT PRIMARY KEY,
	category TEXT NOT NULL DEFAULT 'unspecified'
);
CREATE TABLE IF NOT EXISTS custom_types (
	name TEXT PRIMARY KEY
);

INSERT OR IGNORE INTO config (` + "`key`" + `, value) VALUES
	('compaction_enabled', 'false'),
	('compact_tier1_days', '30'),
	('compact_tier1_dep_levels', '2'),
	('compact_tier2_days', '90'),
	('compact_tier2_dep_levels', '5'),
	('compact_tier2_commits', '100'),
	('compact_batch_size', '50'),
	('compact_parallel_workers', '5'),
	('auto_compact_enabled', 'false');
`
