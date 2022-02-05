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

	px := 0
	sx := 0

	ptmp := -1
	stmp := -1

	if pathname {
		for {
			for charAt(px, pattern) == '*' && charAt(px+1, pattern) == '*' && charAt(px+2, pattern) == '/' {
				px += 3
				ptmp = px
				stmp = sx
			}

			if fnmatchHelper(pattern[px:], str[sx:], flags) {
				for sx < len(str) && charAt(sx, str) != '/' {
					_, ssz := utf8.DecodeRuneInString(str[sx:])
					sx += ssz
				}

				if px < len(pattern) && sx < len(str) {
					px++
					sx++
					continue
				}

				if px >= len(pattern) && sx >= len(pattern) {
					return true
				}
			}

			if ptmp > 0 && stmp > 0 && !(!hasFlag(flags, FNM_DOTMATCH) && charAt(stmp, str) == '.') {
				for stmp < len(str) && charAt(stmp, str) != '/' {
					_, ssz := utf8.DecodeRuneInString(str[stmp:])
					stmp += ssz
				}

				if stmp < len(str) {
					px = ptmp
					stmp++
					sx = stmp
					continue
				}
			}

			return false
		}
	}

	return fnmatchHelper(pattern, str, flags)
}

func fnmatchHelper(pattern string, str string, flags int) bool {
	px := 0
	sx := 0

	ptmp := -1
	stmp := -1

	var prn, srn rune
	var psz, ssz int

	if !hasFlag(flags, FNM_DOTMATCH) && charAt(sx, str) == '.' && charAt(unescape(px, pattern, flags), pattern) != '.' {
		return false
	}

	for {
		switch charAt(px, pattern) {
		case '*':
			for charAt(px, pattern) == '*' {
				px++
			}

			if isEnd(unescape(px, pattern, flags), pattern, flags) {
				px = unescape(px, pattern, flags)
				return true
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
			_, ssz = utf8.DecodeRuneInString(str[sx:])
			sx += ssz
			continue
		case '[':
			if isEnd(sx, str, flags) {
				return false
			}

			if tx, ok := bracketMatch(px+1, pattern, sx, str, flags); ok {
				px = tx
				_, ssz = utf8.DecodeRuneInString(str[sx:])
				sx += ssz
				continue
			}

			goto failed
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
			continue
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
			continue
		}

		return false
	}
}

func bracketMatch(px int, pattern string, sx int, str string, flags int) (int, bool) {
	var ok, not int
	if px >= len(pattern) {
		return 0, false
	}

	if pattern[px] == '!' || pattern[px] == '^' {
		not = 1
		px++
	}

	var t1rn, t2rn rune
	var t1sz, t2sz int
	for pattern[px] != ']' {
		t1 := px
		if !hasFlag(flags, FNM_NOESCAPE) && charAt(t1, pattern) == '\\' {
			t1++
		}

		if t1 >= len(pattern) {
			return 0, false
		}

		t1rn, t1sz = utf8.DecodeRuneInString(pattern[t1:])
		px = t1 + t1sz
		if px >= len(pattern) {
			return 0, false
		}

		if charAt(px, pattern) == '-' && charAt(px+1, pattern) != ']' {
			t2 := px + 1

			if !hasFlag(flags, FNM_NOESCAPE) && charAt(t2, pattern) == '\\' {
				t2++
			}

			if t2 >= len(pattern) {
				return 0, false
			}

			t2rn, t2sz = utf8.DecodeRuneInString(pattern[t2:])
			px = t2 + t2sz
			if ok > 0 {
				continue
			}

			srn, ssz := utf8.DecodeRuneInString(str[sx:])
			if (t1sz == ssz && t1rn == srn) || (t2sz == ssz && t2rn == srn) {
				ok = 1
				continue
			}

			trn, _ := utf8.DecodeRuneInString(pattern[t1:])
			if hasFlag(flags, FNM_CASEFOLD) {
				srn = unicode.ToUpper(srn)
				trn = unicode.ToUpper(trn)
			}
			if srn < trn {
				continue
			}

			trn, _ = utf8.DecodeRuneInString(pattern[t2:])
			if hasFlag(flags, FNM_CASEFOLD) {
				trn = unicode.ToUpper(trn)
			}
			if srn > trn {
				continue
			}

		} else {
			if ok > 0 {
				continue
			}

			srn, ssz := utf8.DecodeRuneInString(str[sx:])
			if t1sz == ssz && t1rn == srn {
				ok = 1
				continue
			}

			if !hasFlag(flags, FNM_CASEFOLD) {
				continue
			}

			prn, _ := utf8.DecodeRuneInString(pattern[px:])
			prn = unicode.ToUpper(prn)
			srn = unicode.ToUpper(srn)

			if prn != srn {
				continue
			}
		}
		ok = 1
	}

	if ok == not {
		return 0, false
	}

	return px + 1, true
}

func hasFlag(mask int, flag int) bool {
	return mask&flag != 0
}

func unescape(p int, pattern string, flags int) int {
	if p >= len(pattern) {
		return p
	}

	if !hasFlag(flags, FNM_NOESCAPE) && pattern[p] == '\\' {
		return p + 1
	}

	return p
}

func isEnd(p int, str string, flags int) bool {
	return p >= len(str) || (hasFlag(flags, FNM_PATHNAME) && str[p] == '/')
}

func charAt(idx int, str string) uint8 {
	if idx >= len(str) {
		return 0
	}

	return str[idx]
}
