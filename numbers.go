// Copyright (C) 2025 arran4
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
	ignoreSpace := false
	underscoreSpace := false
	dotPauses := false
	for _, op := range ops {
		switch op {
		case OpIgnoreSpace:
			ignoreSpace = true
		case OpUnderscoreSpace:
			underscoreSpace = true
		case OpDotPauses:
			dotPauses = true
		}
	}
	p := 0
	var prev rune = 0
	for _, r := range s {
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
			result[p] = r
			p++
			prev = r
			continue
		}
		if prev == rune(s[0]) {
			if dotPauses {
				result[p] = '.'
			} else {
				result[p] = ' '
			}
			p++
		}
		prev = rune(s[0])
		copy(result[p:], []rune(s))
		p += len(s)
	}
	return string(result)
}

func Length(s string) (p int) {
	var prev rune = 0
	for _, r := range s {
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
