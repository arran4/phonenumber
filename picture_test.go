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
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arran4/golang-wordwrap/util"
)

func TestDrawPhoneWithText(t *testing.T) {
	gr, err := util.OpenFont("goregular")
	if err != nil {
		t.Fatalf("failed to open font: %v", err)
	}
	fce := util.GetFontFace(12, 180, gr)

	t.Run("success cases", func(t *testing.T) {
		tempDir := t.TempDir()

		tests := []struct {
			name string
			text string
		}{
			{"empty text", ""},
			{"normal text", "Hello how are you?"},
			{"long text", strings.Repeat("Long text goes here ", 10)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				outPath := filepath.Join(tempDir, strings.ReplaceAll(tt.name, " ", "_")+".png")

				err := DrawPhoneWithText(tt.text, outPath, fce)
				if err != nil {
					t.Fatalf("DrawPhoneWithText failed: %v", err)
				}

				f, err := os.Open(outPath)
				if err != nil {
					t.Fatalf("failed to open generated file: %v", err)
				}
				defer func() {
					_ = f.Close()
				}()

				img, err := png.Decode(f)
				if err != nil {
					t.Fatalf("failed to decode generated PNG: %v", err)
				}

				bounds := img.Bounds()
				if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
					t.Errorf("expected non-zero bounds, got %v", bounds)
				}
			})
		}
	})

	t.Run("invalid output path", func(t *testing.T) {
		tempDir := t.TempDir()
		// Try to write to a directory path which should fail
		outPath := tempDir

		err := DrawPhoneWithText("test", outPath, fce)
		if err == nil {
			t.Errorf("expected error writing to directory path, got nil")
		} else if !strings.Contains(err.Error(), "writing file") {
			t.Errorf("expected error to contain 'writing file', got: %v", err)
		}
	})
}
