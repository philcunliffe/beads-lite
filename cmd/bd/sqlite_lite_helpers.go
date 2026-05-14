//go:build sqlite_lite

package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func maybeAutoBackup(ctx context.Context) {
}

func isDivergedHistoryErr(err error) bool {
	return false
}

func printDivergedHistoryGuidance(operation string) {
}

func atomicWriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}

func joinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

func sanitizeLitePrefix(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "bd"
	}
	return out
}
