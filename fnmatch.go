package fnmatch

const (
	// FNM_NOESCAPE If this flag is set, treat backslash as an ordinary
	// character, instead of an escape character.
	FNM_NOESCAPE = (1 << iota)
	// FNM_PATHNAME If this flag is set, match a slash in string only with a
	// slash in pattern and not by an asterisk (*) or a question
	// mark (?) metacharacter, nor by a bracket expression ([])
	// containing a slash.
	FNM_PATHNAME
	// FNM_PERIOD If this flag is set, a leading period in string has to be
	// matched exactly by a period in pattern.  A period is
	// considered to be leading if it is the first character in
	// string, or if both FNM_PATHNAME is set and the period
	// immediately follows a slash.
	FNM_PERIOD
	// FNM_LEADING_DIR If this flag (a GNU extension) is set, the pattern is
	// considered to be matched if it matches an initial segment
	// of string which is followed by a slash.  This flag is
	// mainly for the internal use of glibc and is implemented
	// only in certain cases.
	FNM_LEADING_DIR
	// FNM_CASEFOLD If this flag (a GNU extension) is set, the pattern is
	// matched case-insensitively.
	FNM_CASEFOLD
	// FNM_EXTGLOB If this flag is set, extended patterns are supported. The
	// extended format is {pattern-list},  with pattern-list being a ','
	// separated list of patterns.
	FNM_EXTGLOB

	FNM_DOTMATCH   = FNM_PERIOD
	FNM_IGNORECASE = FNM_CASEFOLD
	FNM_FILE_NAME  = FNM_PATHNAME
)

// Matches the pattern against the string, with the given flags, and returns true if the match is
// successful.
func Match(pattern string, str string, flags ...int) bool {
	f := clampFlags(flags)
	if hasFlag(f, FNM_EXTGLOB) {
		return expandedMatch(pattern, str, f)
	}
	return match(pattern, str, f)
}
