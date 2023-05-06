package phonenumber

import (
	"strings"
	"unicode"
)

var (
	seq    = "1abc2def3ghi4jkl5mno6pqrs7tuv8wxyz9+* 0"
	lookup = MakeMap(seq)
)

func MakeMap(s string) (result map[rune]string) {
	result = map[rune]string{}
	p := 0
	runes := []rune(s)
	for ri, r := range runes {
		if r == '*' || unicode.IsNumber(r) {
			for rsi, ers := range runes[p : ri+1] {
				result[ers] = strings.Repeat(string([]rune{r}), rsi+1)
			}
			p = ri + 1
		}
	}
	return
}

const OpIgnoreSpace = "IgnoreSpace"
const OpUnderscoreSpace = "UnderscoreSpace"
const OpDotPauses = "DotPauses"

func Numbers(s string, ops ...any) string {
	result := []rune(strings.Repeat(" ", Length(s)))
	runes := []rune(s)
	ignoreSpace := false
	underscoreSpace := false
	dotPauses := false
	for _, op := range ops {
		if op == OpIgnoreSpace {
			ignoreSpace = true
		} else if op == OpUnderscoreSpace {
			underscoreSpace = true
		} else if op == OpDotPauses {
			dotPauses = true
		}
	}
	p := 0
	var prev rune = 0
	for _, r := range runes {
		if r == ' ' {
			if ignoreSpace {
				p++
				prev = r
				continue
			} else if underscoreSpace {
				result[p] = '_'
				p++
				prev = r
				continue
			}
		}
		lcr := unicode.ToLower(r)
		s, ok := lookup[lcr]
		if !ok {
			p++
			prev = r
			continue
		}
		if prev == rune(s[0]) {
			p++
			if dotPauses {
				result[p] = '.'
			} else {
				result[p] = ' '
			}
		}
		prev = rune(s[0])
		copy(result[p:], []rune(s))
		p += len(s)
	}
	return string(result)
}

func Length(s string) (p int) {
	runes := []rune(s)
	var prev rune = 0
	for _, r := range runes {
		lcr := unicode.ToLower(r)
		s, ok := lookup[lcr]
		if !ok {
			p++
			prev = r
			continue
		}
		p += len(s)
		if prev == rune(s[0]) {
			p++
		}
		prev = rune(s[0])
	}
	return
}
