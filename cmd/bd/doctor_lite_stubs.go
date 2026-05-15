package main

import "fmt"

// previewFixes is the lite-mode replacement for the legacy Dolt fix preview.
// It surfaces the doctor result without attempting any repair.
func previewFixes(result doctorResult) {
	if len(result.Checks) == 0 {
		return
	}
	fmt.Println("Doctor --fix is not implemented in lite mode; no repairs will be applied.")
}

// applyFixes is the lite-mode replacement for the legacy Dolt --fix path.
func applyFixes(result doctorResult) {
	previewFixes(result)
}

// trackBdVersion was the legacy version bookkeeping helper. In lite mode the
// per-DB version table is updated lazily on first write; this stub keeps the
// doctor command flow intact.
func trackBdVersion() {}

// autoMigrateOnVersionBump used to trigger schema migrations when the binary
// version moved forward. SQLite migrations run on store open in the lite build.
func autoMigrateOnVersionBump(_ string) {}

// runCheckHealth is the lite-mode --check-health hook. The lite build does not
// run any external server, so this is a quiet success.
func runCheckHealth(_ string) {}

// runDeepValidation is the lite-mode --deep hook.
func runDeepValidation(_ string) {
	fmt.Println("Deep validation is not available in lite mode.")
}

// runServerHealth is the lite-mode --server hook. SQLite-only, no server.
func runServerHealth(_ string) {
	fmt.Println("Server-mode health checks are not available in lite mode.")
}
