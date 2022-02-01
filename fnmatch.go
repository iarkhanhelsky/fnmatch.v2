package fnmatch

const (
	FNM_NOESCAPE = (1 << iota)
	FNM_PATHNAME
	FNM_PERIOD

	FNM_LEADING_DIR
	FNM_CASEFOLD

	FNM_IGNORECASE = FNM_CASEFOLD
	FNM_FILE_NAME  = FNM_PATHNAME
)

// Matches the pattern against the string, with the given flags, and returns true if the match is
// successful.
func Match(pattern, s string, flags int) bool {
	return false
}
