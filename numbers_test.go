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
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestLength(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		wantP int
	}{
		{
			name:  "A is 1",
			s:     "a",
			wantP: 1,
		},
		{
			name:  "b is 2",
			s:     "b",
			wantP: 2,
		},
		{
			name:  "space is 1",
			s:     " ",
			wantP: 1,
		},
		{
			name:  "nw is 2",
			s:     "mw",
			wantP: 2,
		},
		{
			name:  "Empty is 0",
			s:     "",
			wantP: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotP := Length(tt.s); gotP != tt.wantP {
				t.Errorf("Length() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

func TestMakeMap(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		wantResult map[rune]string
	}{
		{
			name: "2s",
			s:    "abc2",
			wantResult: map[rune]string{
				'a': "2",
				'b': "22",
				'c': "222",
				'2': "2222",
			},
		},
		{
			name: "2s 3s",
			s:    "abc2def3",
			wantResult: map[rune]string{
				'a': "2",
				'b': "22",
				'c': "222",
				'2': "2222",
				'd': "3",
				'e': "33",
				'f': "333",
				'3': "3333",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := MakeMap(tt.s)
			if d := cmp.Diff(tt.wantResult, gotResult); d != "" {
				t.Errorf("MakeMap() = %v, want %v: \n%s", gotResult, tt.wantResult, d)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
		args []Option
	}{
		{
			name: "abc",
			s:    "abc",
			want: "2 22 222",
		},
		{
			name: "Uppercase ABC",
			s:    "ABC",
			want: "2 22 222",
		},
		{
			name: "aei",
			s:    "aei",
			want: "233444",
		},
		{
			name: "123",
			s:    "123",
			want: "122223333",
		},
		{
			name: "Spaces",
			s:    "a d gj",
			want: "203045",
		},
		{
			name: "Spaces Ignore spaces",
			s:    "a d gj",
			want: "2 3 45",
			args: []Option{WithIgnoreSpace()},
		},
		{
			name: "Underscore spaces",
			s:    "a d gj",
			want: "2_3_45",
			args: []Option{WithUnderscoreSpace()},
		},
		{
			name: "Dot pauses",
			s:    "hello",
			want: "4433555.555666",
			args: []Option{WithDotPauses()},
		},
		{
			name: "Pass through unknowns",
			s:    "hello?",
			want: "4433555 555666?",
		},
		{
			name: "Unicode outside mapping",
			s:    "aπb",
			want: "2π22",
		},
		{
			name: "+ and *",
			s:    "+*",
			want: "* **",
		},
		{
			name: "Empty string",
			s:    "",
			want: "",
		},
		{
			name: "All options",
			s:    "hello world",
			want: "4433555.555666 96667775553",
			args: []Option{WithIgnoreSpace(), WithUnderscoreSpace(), WithDotPauses()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Convert(tt.s, tt.args...); got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumbers(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
		args []any
	}{
		{
			name: "abc",
			s:    "abc",
			want: "2 22 222",
		},
		{
			name: "Spaces Ignore spaces",
			s:    "a d gj",
			want: "2 3 45",
			args: []any{OpIgnoreSpace},
		},
		{
			name: "Dot pauses",
			s:    "hello",
			want: "4433555.555666",
			args: []any{OpDotPauses},
		},
		{
			name: "Underscore space",
			s:    "a b",
			want: "2_22",
			args: []any{OpUnderscoreSpace},
		},
		{
			name: "Unknown options ignored",
			s:    "a b",
			want: "2022",
			args: []any{"SomeUnknownOption", 123},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Numbers(tt.s, tt.args...); got != tt.want {
				t.Errorf("Numbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func FuzzConvertAndNumbers(f *testing.F) {
	f.Add("hello world")
	f.Add("abc 123 !@#")
	f.Add("")
	f.Add("A π C")

	f.Fuzz(func(t *testing.T, s string) {
		// Test no panics for arbitrary UTF-8
		_ = Convert(s)
		_ = Convert(s, WithIgnoreSpace())
		_ = Convert(s, WithUnderscoreSpace())
		_ = Convert(s, WithDotPauses())

		// Test equivalence between legacy Numbers and new Convert
		if got, want := Numbers(s), Convert(s); got != want {
			t.Errorf("Numbers(%q) != Convert(%q): %v != %v", s, s, got, want)
		}
		if got, want := Numbers(s, OpIgnoreSpace), Convert(s, WithIgnoreSpace()); got != want {
			t.Errorf("Numbers(%q, OpIgnoreSpace) != Convert(...): %v != %v", s, got, want)
		}
		if got, want := Numbers(s, OpUnderscoreSpace), Convert(s, WithUnderscoreSpace()); got != want {
			t.Errorf("Numbers(%q, OpUnderscoreSpace) != Convert(...): %v != %v", s, got, want)
		}
		if got, want := Numbers(s, OpDotPauses), Convert(s, WithDotPauses()); got != want {
			t.Errorf("Numbers(%q, OpDotPauses) != Convert(...): %v != %v", s, got, want)
		}
		if got, want := Numbers(s, OpIgnoreSpace, OpDotPauses), Convert(s, WithIgnoreSpace(), WithDotPauses()); got != want {
			t.Errorf("Numbers(...) != Convert(...): %v != %v", got, want)
		}
	})
}
