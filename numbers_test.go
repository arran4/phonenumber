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
			name:  "b is 1",
			s:     "b",
			wantP: 1,
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

func TestNumbers(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Numbers(tt.s); got != tt.want {
				t.Errorf("Numbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
