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
	"path/filepath"
	"testing"
)

func TestAppModelValidation(t *testing.T) {
	tests := []struct {
		name         string
		fontSizeText string
		filename     string
		expectError  bool
	}{
		{
			name:         "Valid model",
			fontSizeText: "12",
			filename:     "out.png",
			expectError:  false,
		},
		{
			name:         "Empty font size",
			fontSizeText: "",
			filename:     "out.png",
			expectError:  true,
		},
		{
			name:         "Invalid font size format",
			fontSizeText: "abc",
			filename:     "out.png",
			expectError:  true,
		},
		{
			name:         "Zero font size",
			fontSizeText: "0",
			filename:     "out.png",
			expectError:  true,
		},
		{
			name:         "Negative font size",
			fontSizeText: "-5",
			filename:     "out.png",
			expectError:  true,
		},
		{
			name:         "Empty filename",
			fontSizeText: "12",
			filename:     "",
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewAppModel()
			m.UpdateFontSize(tc.fontSizeText)
			m.UpdateFilename(tc.filename)

			err := m.Validate()
			if tc.expectError && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestAppModelUpdateText(t *testing.T) {
	m := NewAppModel()
	m.UpdateText("hello")

	if m.Text != "hello" {
		t.Errorf("Expected Text to be 'hello', got %s", m.Text)
	}
	if m.Result == "" {
		t.Errorf("Expected Result to not be empty")
	}
}

func TestAppModelGenerate(t *testing.T) {
	m := NewAppModel()
	m.UpdateText("test")
	m.UpdateFontSize("12")
	m.UpdateFilename(filepath.Join(t.TempDir(), "test.png"))

	err := m.Generate()
	if err != nil {
		t.Errorf("Expected Generate to succeed, but got error: %v", err)
	}
	if m.Error != nil {
		t.Errorf("Expected model Error to be nil, got: %v", m.Error)
	}
	if m.Status == "" {
		t.Errorf("Expected model Status to be populated on success")
	}
}

func TestAppModelGenerateErrorPropagation(t *testing.T) {
	m := NewAppModel()
	m.UpdateText("test")
	m.UpdateFontSize("-1")

	err := m.Generate()
	if err == nil {
		t.Errorf("Expected Generate to fail due to invalid font size")
	}
	if m.Error == nil {
		t.Errorf("Expected model Error to be populated")
	}
}
