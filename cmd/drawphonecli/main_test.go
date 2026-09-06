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

package main

import (
	"bytes"
	"flag"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		args        []string
		expectedErr error
		errContains string
		outContains string
	}{
		{
			name:        "help flag",
			args:        []string{"-h"},
			expectedErr: flag.ErrHelp,
		},
		{
			name:        "invalid flag",
			args:        []string{"-invalidflag"},
			errContains: "flag provided but not defined: -invalidflag",
		},
		{
			name:        "invalid fontsize",
			args:        []string{"-fontsize", "invalid"},
			errContains: "invalid value \"invalid\" for flag -fontsize",
		},
		{
			name:        "invalid output dir",
			args:        []string{"-outfile", tempDir},
			errContains: "drawing phone: writing file",
		},
		{
			name:        "default success",
			args:        []string{"-outfile", filepath.Join(tempDir, "default.png")},
			outContains: "Wrote: " + filepath.Join(tempDir, "default.png"),
		},
		{
			name:        "specific text",
			args:        []string{"-outfile", filepath.Join(tempDir, "specific.png"), "-text", "abc"},
			outContains: "'2.22.222'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			err := run(tt.args, &stdout, &stderr)

			if tt.expectedErr != nil {
				if err != tt.expectedErr {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if tt.errContains != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got: %v", tt.errContains, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.outContains != "" {
				out := stdout.String()
				if !strings.Contains(out, tt.outContains) {
					t.Errorf("expected stdout to contain %q, got: %s", tt.outContains, out)
				}
			}
		})
	}
}
