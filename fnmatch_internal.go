package fnmatch

import (
	"unicode"
	"unicode/utf8"
)

func match(pattern string, str string, flags int) bool {
	px := 0
	sx := 0

	ptmp := -1
	stmp := -1

	if hasFlag(flags, FNM_PATHNAME) {
		for {
			for charAt(px, pattern) == '*' && charAt(px+1, pattern) == '*' && charAt(px+2, pattern) == '/' {
				px += 3
				ptmp = px
				stmp = sx
			}

			ok, tp, ts := fnmatchHelper(pattern[px:], str[sx:], flags)
			px += tp
			sx += ts

			if ok {
				for sx < len(str) && charAt(sx, str) != '/' {
					_, ssz := utf8.DecodeRuneInString(str[sx:])
					sx += ssz
				}

				if px < len(pattern) && sx < len(str) {
					px++
					sx++
					continue
				}

				if px >= len(pattern) && sx >= len(str) {
					return true
				}
			}

			if ptmp >= 0 && stmp >= 0 && !(!hasFlag(flags, FNM_DOTMATCH) && charAt(stmp, str) == '.') {
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

	ok, _, _ := fnmatchHelper(pattern, str, flags)
	return ok
}

func fnmatchHelper(pattern string, str string, flags int) (bool, int, int) {
	px := 0
	sx := 0

	ptmp := -1
	stmp := -1

	var prn, srn rune
	var psz, ssz int

	if !hasFlag(flags, FNM_DOTMATCH) && charAt(sx, str) == '.' && charAt(unescape(px, pattern, flags), pattern) != '.' {
		return false, px, sx
	}

	for {
		switch charAt(px, pattern) {
		case '*':
			for charAt(px, pattern) == '*' {
				px++
			}

			if isEnd(unescape(px, pattern, flags), pattern, flags) {
				px = unescape(px, pattern, flags)
				return true, px, sx
			}

			if isEnd(sx, str, flags) {
				return false, px, sx
			}

			ptmp = px
			stmp = sx
		case '?':
			if isEnd(sx, str, flags) {
				return false, px, sx
			}

			px++
			_, ssz = utf8.DecodeRuneInString(str[sx:])
			sx += ssz
			continue
		case '[':
			if isEnd(sx, str, flags) {
				return false, px, sx
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
			return isEnd(px, pattern, flags), px, sx
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
			_, ssz = utf8.DecodeRuneInString(str[sx:])
			sx += ssz
			continue
		}

		return false, px, sx
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

	return px + 1, ok != not
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

func clampFlags(flags []int) int {
	o := 0
	for _, f := range flags {
		o |= f
	}

	return o
}
