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
	"image/png"
	"os"
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
		{
			name:        "custom fontsize",
			args:        []string{"-outfile", filepath.Join(tempDir, "fontsize.png"), "-fontsize", "24"},
			outContains: "Wrote:",
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

func TestExecuteExitStatus(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		args         []string
		expectedCode int
	}{
		{
			name:         "help flag returns 0",
			args:         []string{"-h"},
			expectedCode: 0,
		},
		{
			name:         "invalid flag returns 1",
			args:         []string{"-invalidflag"},
			expectedCode: 1,
		},
		{
			name:         "valid args return 0",
			args:         []string{"-outfile", filepath.Join(tempDir, "execute_success.png")},
			expectedCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock os.Args
			importArgs := append([]string{"drawphonecli"}, tt.args...)

			// We need to temporarily replace os.Args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = importArgs

			// Capture stdout and stderr to avoid spamming the test output
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			defer func() {
				os.Stdout = oldStdout
				os.Stderr = oldStderr
			}()

			code := execute()
			if code != tt.expectedCode {
				t.Errorf("expected exit code %d, got %d", tt.expectedCode, code)
			}
		})
	}
}

func TestFontSizeStructuralDifference(t *testing.T) {
	tempDir := t.TempDir()

	outSmall := filepath.Join(tempDir, "small.png")
	outLarge := filepath.Join(tempDir, "large.png")

	// Render small text
	var stdout1, stderr1 bytes.Buffer
	err := run([]string{"-outfile", outSmall, "-fontsize", "10", "-text", "Hello"}, &stdout1, &stderr1)
	if err != nil {
		t.Fatalf("failed to run small text: %v", err)
	}

	// Render large text
	var stdout2, stderr2 bytes.Buffer
	err = run([]string{"-outfile", outLarge, "-fontsize", "100", "-text", "Hello"}, &stdout2, &stderr2)
	if err != nil {
		t.Fatalf("failed to run large text: %v", err)
	}

	fSmall, err := os.Open(outSmall)
	if err != nil {
		t.Fatalf("failed to open small png: %v", err)
	}
	defer func() { _ = fSmall.Close() }()
	imgSmall, err := png.Decode(fSmall)
	if err != nil {
		t.Fatalf("failed to decode small png: %v", err)
	}

	fLarge, err := os.Open(outLarge)
	if err != nil {
		t.Fatalf("failed to open large png: %v", err)
	}
	defer func() { _ = fLarge.Close() }()
	imgLarge, err := png.Decode(fLarge)
	if err != nil {
		t.Fatalf("failed to decode large png: %v", err)
	}

	// Since "Hello" with size 10 vs 100 might word-wrap differently:
	// A much larger font size will either be taller overall, or significantly wider overall/scale text area.
	// We'll check the text drawing area bounds difference if they are different heights.
	// Actually word-wrapping means large font might take multiple lines, so it should be taller
	if imgLarge.Bounds().Dy() == imgSmall.Bounds().Dy() {
		t.Errorf("expected font sizes to result in different image dimensions (height). Both are: %d", imgLarge.Bounds().Dy())
	}
}
