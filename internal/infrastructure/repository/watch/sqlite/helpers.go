package sqlite

// boolToInt64 converts bool to int64 for SQLite
// stores booleans as INTEGER (0 or 1)
func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
