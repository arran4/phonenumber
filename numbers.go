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

// Package phonenumber provides functionality for generating and drawing Nokia-style phone number key sequences.
package phonenumber

import (
	"strings"
	"unicode"
)

var (
	seq    = "1abc2def3ghi4jkl5mno6pqrs7tuv8wxyz9+* 0"
	lookup = MakeMap(seq)
)

// MakeMap creates a lookup map from a character sequence definition.
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

// options holds the configuration for conversion.
type options struct {
	ignoreSpace     bool
	underscoreSpace bool
	dotPauses       bool
}

// Option defines a type-safe functional option for configuring text conversion.
type Option func(*options)

// WithIgnoreSpace returns an option to ignore spaces when translating.
func WithIgnoreSpace() Option {
	return func(o *options) {
		o.ignoreSpace = true
	}
}

// WithUnderscoreSpace returns an option to use underscores for spaces.
func WithUnderscoreSpace() Option {
	return func(o *options) {
		o.underscoreSpace = true
	}
}

// WithDotPauses returns an option to use dots for pauses.
func WithDotPauses() Option {
	return func(o *options) {
		o.dotPauses = true
	}
}

// OpIgnoreSpace defines an option to ignore spaces when translating.
//
// Deprecated: Use WithIgnoreSpace() with Convert instead.
const OpIgnoreSpace = "IgnoreSpace"

// OpUnderscoreSpace defines an option to use underscores for spaces.
//
// Deprecated: Use WithUnderscoreSpace() with Convert instead.
const OpUnderscoreSpace = "UnderscoreSpace"

// OpDotPauses defines an option to use dots for pauses.
//
// Deprecated: Use WithDotPauses() with Convert instead.
const OpDotPauses = "DotPauses"

// Convert translates a string into its corresponding numeric key sequence using a type-safe API.
//
// Default Behaviour:
// - Space characters (' ') are mapped to '0' (since the Nokia '0' key represents space).
// - Unknown punctuation or characters not present in the Nokia keypad are passed through as literal runes.
//
// Option Interactions & Precedence:
// - WithIgnoreSpace(): Ignores the default mapping to '0' and retains literal spaces (' '). Takes highest precedence over WithUnderscoreSpace().
// - WithUnderscoreSpace(): Replaces spaces with underscores ('_') instead of mapping to '0'.
//   - WithDotPauses(): When identical digits are repeated sequentially due to adjacent characters mapping
//     to the same keypad key, a pause is inserted between them. By default, this pause is a space (' ').
//     WithDotPauses() alters this pause character to a dot ('.').
func Convert(s string, opts ...Option) string {
	var config options
	for _, opt := range opts {
		opt(&config)
	}
	result := []rune(strings.Repeat(" ", Length(s)))
	p := 0
	var prev rune
	for _, r := range s {
		if r == ' ' {
			if config.ignoreSpace {
				p++
				prev = r
				continue
			} else if config.underscoreSpace {
				result[p] = '_'
				p++
				prev = r
				continue
			}
		}
		lcr := unicode.ToLower(r)
		lookupStr, ok := lookup[lcr]
		if !ok {
			result[p] = r
			p++
			prev = r
			continue
		}
		if prev == rune(lookupStr[0]) {
			if config.dotPauses {
				result[p] = '.'
			} else {
				result[p] = ' '
			}
			p++
		}
		prev = rune(lookupStr[0])
		copy(result[p:], []rune(lookupStr))
		p += len(lookupStr)
	}
	return string(result)
}

// Numbers translates a string into its corresponding numeric key sequence.
//
// Deprecated: Use Convert with type-safe Option builders instead.
func Numbers(s string, ops ...any) string {
	var opts []Option
	for _, op := range ops {
		switch op {
		case OpIgnoreSpace:
			opts = append(opts, WithIgnoreSpace())
		case OpUnderscoreSpace:
			opts = append(opts, WithUnderscoreSpace())
		case OpDotPauses:
			opts = append(opts, WithDotPauses())
		}
	}
	return Convert(s, opts...)
}

// Length calculates the required length of the translated string array.
func Length(s string) (p int) {
	var prev rune
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
