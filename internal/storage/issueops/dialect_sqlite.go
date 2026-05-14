//go:build sqlite_lite

package issueops

import "fmt"

func insertIssueSQL(table string) string {
	return fmt.Sprintf(`
		INSERT INTO %s (
			id, content_hash, title, description, design, acceptance_criteria, notes,
			status, priority, issue_type, assignee, estimated_minutes,
			created_at, created_by, owner, updated_at, started_at, closed_at, external_ref, spec_id,
			compaction_level, compacted_at, compacted_at_commit, original_size,
			sender, ephemeral, no_history, wisp_type, pinned, is_template,
			mol_type, work_type, source_system, source_repo, close_reason,
			event_kind, actor, target, payload,
			await_type, await_id, timeout_ns, waiters,
			due_at, defer_until, metadata
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?
		)
		ON CONFLICT(id) DO UPDATE SET
			content_hash = excluded.content_hash,
			title = excluded.title,
			description = excluded.description,
			design = excluded.design,
			acceptance_criteria = excluded.acceptance_criteria,
			notes = excluded.notes,
			status = excluded.status,
			priority = excluded.priority,
			issue_type = excluded.issue_type,
			assignee = excluded.assignee,
			estimated_minutes = excluded.estimated_minutes,
			updated_at = excluded.updated_at,
			started_at = excluded.started_at,
			closed_at = excluded.closed_at,
			external_ref = excluded.external_ref,
			source_repo = excluded.source_repo,
			close_reason = excluded.close_reason,
			metadata = excluded.metadata
	`, table)
}

func insertLabelSQL(table string) string {
	return fmt.Sprintf(`INSERT OR IGNORE INTO %s (issue_id, label) VALUES (?, ?)`, table)
}

func upsertLabelSQL(table string) string {
	return fmt.Sprintf(`
		INSERT INTO %s (issue_id, label)
		VALUES (?, ?)
		ON CONFLICT(issue_id, label) DO NOTHING
	`, table)
}

func insertSelectIgnoreSQL(table, columns, selectSQL string) string {
	return fmt.Sprintf("INSERT OR IGNORE INTO %s (%s) %s", table, columns, selectSQL)
}

func upsertDependencyNoopSQL(table string) string {
	return fmt.Sprintf(`
		INSERT INTO %s (issue_id, depends_on_id, type, created_by, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(issue_id, depends_on_id) DO NOTHING
	`, table)
}

func upsertChildCounterMaxSQL() string {
	return `
		INSERT INTO child_counters (parent_id, last_child) VALUES (?, ?)
		ON CONFLICT(parent_id) DO UPDATE SET last_child = MAX(child_counters.last_child, excluded.last_child)
	`
}

func setChildCounterSQL() string {
	return `
		INSERT INTO child_counters (parent_id, last_child) VALUES (?, ?)
		ON CONFLICT(parent_id) DO UPDATE SET last_child = excluded.last_child
	`
}

func insertDependencyNowSQL(table string) string {
	return fmt.Sprintf(`
		INSERT INTO %s (issue_id, depends_on_id, type, created_at, created_by, metadata, thread_id)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?, ?, ?)
	`, table)
}

func upsertRepoMtimeSQL() string {
	return `
		INSERT INTO repo_mtimes (repo_path, jsonl_path, mtime_ns, last_checked)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(repo_path) DO UPDATE SET
			jsonl_path = excluded.jsonl_path,
			mtime_ns = excluded.mtime_ns,
			last_checked = CURRENT_TIMESTAMP
	`
}

func upsertFederationPeerSQL() string {
	return `
		INSERT INTO federation_peers (name, remote_url, username, password_encrypted, sovereignty)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			remote_url = excluded.remote_url,
			username = excluded.username,
			password_encrypted = excluded.password_encrypted,
			sovereignty = excluded.sovereignty,
			updated_at = CURRENT_TIMESTAMP
	`
}

func insertDefaultConfigSQL() string {
	return `
		INSERT OR IGNORE INTO config (` + "`key`" + `, value) VALUES
			('compaction_enabled', 'false'),
			('compact_tier1_days', '30'),
			('compact_tier1_dep_levels', '2'),
			('compact_tier2_days', '90'),
			('compact_tier2_dep_levels', '5'),
			('compact_tier2_commits', '100'),
			('compact_batch_size', '50'),
			('compact_parallel_workers', '5'),
			('auto_compact_enabled', 'false')
	`
}

func jsonPathExistsClause() string {
	return "json_extract(metadata, ?) IS NOT NULL"
}

func jsonPathEqualsClause() string {
	return "json_extract(metadata, ?) = ?"
}

func recentIssueSortClause() string {
	return `
		CASE WHEN created_at >= datetime('now', '-48 hours') THEN 0 ELSE 1 END ASC,
		CASE WHEN created_at >= datetime('now', '-48 hours') THEN priority ELSE 999 END ASC,
	`
}

func currentTimestampSQL() string {
	return "CURRENT_TIMESTAMP"
}

func descendantPathStartExpr() string {
	return "',' || ? || ',' || issue_id || ','"
}

func descendantPathAppendExpr() string {
	return "d.path || e.issue_id || ','"
}

func descendantPathContainsExpr() string {
	return "instr(d.path, ',' || e.issue_id || ',') = 0"
}

func updateDependencyPrefixSQL(column string) string {
	return fmt.Sprintf(`
		UPDATE dependencies
		SET %s = ? || substr(%s, length(?) + 1)
		WHERE %s LIKE ?
	`, column, column, column)
}
