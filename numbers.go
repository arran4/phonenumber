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
			p = ri
		}
	}
	return
}

const OpIgnoreSpace = "IgnoreSpace"

func Numbers(s string, ops ...any) string {
	result := strings.Repeat(" ", len(s)+Length(s))
	runes := []rune(s)
	ignoreSpace := false
	for _, op := range ops {
		if op == OpIgnoreSpace {
			ignoreSpace = true
		}
	}
	p := 0
	for _, r := range runes {
		if r == ' ' && ignoreSpace {
			p++
			continue
		}
		lcr := unicode.ToLower(r)
		s, ok := lookup[lcr]
		if !ok {
			p++
			continue
		}
		copy([]rune(result[p:]), []rune(s))
		p += len(s)
	}
	return result
}

func Length(s string) (p int) {
	runes := []rune(s)
	for _, r := range runes {
		lcr := unicode.ToLower(r)
		s, ok := lookup[lcr]
		if !ok {
			p++
			continue
		}
		p += len(s)
	}
	return
}
