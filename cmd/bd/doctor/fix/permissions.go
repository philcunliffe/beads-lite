package fix

// Permissions repairs filesystem permissions on the beads directory. The lite
// build keeps the entry point so the doctor surface compiles, but it does
// not change anything on disk — pure-Go SQLite has no special permission
// requirements beyond the read/write that os.Create already enforces.
func Permissions(_ string) error { return nil }
