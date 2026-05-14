//go:build !sqlite_lite

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
		ON DUPLICATE KEY UPDATE
			content_hash = VALUES(content_hash),
			title = VALUES(title),
			description = VALUES(description),
			design = VALUES(design),
			acceptance_criteria = VALUES(acceptance_criteria),
			notes = VALUES(notes),
			status = VALUES(status),
			priority = VALUES(priority),
			issue_type = VALUES(issue_type),
			assignee = VALUES(assignee),
			estimated_minutes = VALUES(estimated_minutes),
			updated_at = VALUES(updated_at),
			started_at = VALUES(started_at),
			closed_at = VALUES(closed_at),
			external_ref = VALUES(external_ref),
			source_repo = VALUES(source_repo),
			close_reason = VALUES(close_reason),
			metadata = VALUES(metadata)
	`, table)
}

func insertLabelSQL(table string) string {
	return fmt.Sprintf(`INSERT IGNORE INTO %s (issue_id, label) VALUES (?, ?)`, table)
}

func upsertLabelSQL(table string) string {
	return fmt.Sprintf(`
		INSERT INTO %s (issue_id, label)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE label = label
	`, table)
}

func insertSelectIgnoreSQL(table, columns, selectSQL string) string {
	return fmt.Sprintf("INSERT IGNORE INTO %s (%s) %s", table, columns, selectSQL)
}

func upsertDependencyNoopSQL(table string) string {
	return fmt.Sprintf(`
		INSERT INTO %s (issue_id, depends_on_id, type, created_by, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE type = type
	`, table)
}

func upsertChildCounterMaxSQL() string {
	return `
		INSERT INTO child_counters (parent_id, last_child) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE last_child = GREATEST(last_child, ?)
	`
}

func setChildCounterSQL() string {
	return `
		INSERT INTO child_counters (parent_id, last_child) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE last_child = ?
	`
}

func insertDependencyNowSQL(table string) string {
	return fmt.Sprintf(`
		INSERT INTO %s (issue_id, depends_on_id, type, created_at, created_by, metadata, thread_id)
		VALUES (?, ?, ?, NOW(), ?, ?, ?)
	`, table)
}

func upsertRepoMtimeSQL() string {
	return `
		INSERT INTO repo_mtimes (repo_path, jsonl_path, mtime_ns, last_checked)
		VALUES (?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			jsonl_path = VALUES(jsonl_path),
			mtime_ns = VALUES(mtime_ns),
			last_checked = NOW()
	`
}

func upsertFederationPeerSQL() string {
	return `
		INSERT INTO federation_peers (name, remote_url, username, password_encrypted, sovereignty)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			remote_url = VALUES(remote_url),
			username = VALUES(username),
			password_encrypted = VALUES(password_encrypted),
			sovereignty = VALUES(sovereignty),
			updated_at = CURRENT_TIMESTAMP
	`
}

func insertDefaultConfigSQL() string {
	return `
		INSERT IGNORE INTO config (` + "`key`" + `, value) VALUES
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
	return "JSON_EXTRACT(metadata, ?) IS NOT NULL"
}

func jsonPathEqualsClause() string {
	return "JSON_UNQUOTE(JSON_EXTRACT(metadata, ?)) = ?"
}

func recentIssueSortClause() string {
	return `
		CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 48 HOUR) THEN 0 ELSE 1 END ASC,
		CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 48 HOUR) THEN priority ELSE 999 END ASC,
	`
}

func currentTimestampSQL() string {
	return "UTC_TIMESTAMP()"
}

func descendantPathStartExpr() string {
	return "CONCAT(',', ?, ',', issue_id, ',')"
}

func descendantPathAppendExpr() string {
	return "CONCAT(d.path, e.issue_id, ',')"
}

func descendantPathContainsExpr() string {
	return "LOCATE(CONCAT(',', e.issue_id, ','), d.path) = 0"
}

func updateDependencyPrefixSQL(column string) string {
	return fmt.Sprintf(`
		UPDATE dependencies
		SET %s = CONCAT(?, SUBSTRING(%s, LENGTH(?) + 1))
		WHERE %s LIKE ?
	`, column, column, column)
}
