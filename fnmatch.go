package fnmatch

import (
	"unicode"
	"unicode/utf8"
)

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
	// FNM_EXTMATCH
	// If this flag (a GNU extension) is set, extended patterns
	// are supported, as introduced by 'ksh' and now supported by
	// other shells.  The extended format is as follows, with
	// pattern-list being a '|' separated list of patterns.
	//
	//	'?(pattern-list)'
	//		The pattern matches if zero or one occurrences of any of
	//		the patterns in the pattern-list match the input string.
	//
	//  '*(pattern-list)'
	//		The pattern matches if zero or more occurrences of any of
	//		the patterns in the pattern-list match the input string.
	//	'+(pattern-list)'
	//		The pattern matches if one or more occurrences of any of
	//		the patterns in the pattern-list match the input string.
	//
	// 	'@(pattern-list)'
	//		The pattern matches if exactly one occurrence of any of
	//		the patterns in the pattern-list match the input string.
	//
	//	'!(pattern-list)'
	//		The pattern matches if the input string cannot be matched
	//		with any of the patterns in the pattern-list.
	FNM_EXTMATCH

	FNM_DOTMATCH   = FNM_PERIOD
	FNM_IGNORECASE = FNM_CASEFOLD
	FNM_FILE_NAME  = FNM_PATHNAME
)

// Matches the pattern against the string, with the given flags, and returns true if the match is
// successful.
func Match(pattern string, str string, flags int) bool {
	pathname := hasFlag(flags, FNM_PATHNAME)

	if pathname {
		// TODO
		return false
	}

	return fnmatchHelper(pattern, str, flags)
}

func fnmatchHelper(pattern string, str string, flags int) bool {
	period := !hasFlag(flags, FNM_DOTMATCH)

	px := 0
	sx := 0

	ptmp := -1
	stmp := -1

	var prn, srn rune
	var psz, ssz int

	if period && str[sx] == '.' && pattern[unescape(px, pattern, flags)] != '.' {
		return false
	}

	for {
		switch pattern[px] {
		case '*':
			for pattern[px] == '*' {
				px++
			}

			if isEnd(unescape(px, pattern, flags), pattern, flags) {
				px = unescape(px, pattern, flags)
				return false
			}

			if isEnd(sx, str, flags) {
				return false
			}

			ptmp = px
			stmp = sx
		case '?':
			if isEnd(sx, str, flags) {
				return false
			}

			px++
			_, ssz = utf8.DecodeRuneInString(str[stmp:])
			sx += ssz
			continue
		}

		px = unescape(px, pattern, flags)
		if isEnd(sx, str, flags) {
			return isEnd(px, pattern, flags)
		}
		if isEnd(px, pattern, flags) {
			goto failed
		}

		prn, psz = utf8.DecodeRuneInString(pattern[px:])
		srn, ssz = utf8.DecodeRuneInString(str[sx:])
		if psz == ssz && prn == srn {
			px += psz
			sx += ssz
		}

		if !hasFlag(flags, FNM_CASEFOLD) {
			goto failed
		}

		if unicode.ToUpper(prn) != unicode.ToUpper(srn) {
			goto failed
		}
		px += psz
		sx += ssz
		continue

	failed:
		// try next '*' position
		if ptmp >= 0 && stmp >= 0 {
			px = ptmp
			_, ssz = utf8.DecodeRuneInString(str[stmp:])
			sx += ssz
		}

		return false
	}
}

func hasFlag(mask int, flag int) bool {
	return mask&flag != 0
}

func unescape(p int, pattern string, flags int) int {
	if !hasFlag(flags, FNM_NOESCAPE) && pattern[p] == '\\' {
		return p + 1
	}

	return p
}

func isEnd(p int, pattern string, flags int) bool {
	return p >= len(pattern) || (hasFlag(flags, FNM_PATHNAME) && pattern[p] == '/')
}
