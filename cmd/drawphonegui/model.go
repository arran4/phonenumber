//go:build !headless

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
	"errors"
	"fmt"
	"phonenumber"
	"strconv"

	"github.com/arran4/golang-wordwrap/util"
)

// AppModel holds the state and business logic of the application, completely isolated from GUI toolkit references.
type AppModel struct {
	Text         string
	Result       string
	FontSizeText string
	Filename     string
	Status       string
	Error        error
}

// NewAppModel creates a new AppModel initialized with default values.
func NewAppModel() *AppModel {
	m := &AppModel{
		Text:         "Hello how are you?",
		FontSizeText: "12",
		Filename:     "out.png",
	}
	m.UpdateText(m.Text)
	return m
}

// UpdateText updates the input text and updates the translation result.
func (m *AppModel) UpdateText(text string) {
	m.Text = text
	m.Result = phonenumber.Convert(text, phonenumber.WithIgnoreSpace(), phonenumber.WithDotPauses())
	m.ClearStatus()
}

// UpdateFontSize updates the font size string.
func (m *AppModel) UpdateFontSize(fontSize string) {
	m.FontSizeText = fontSize
	m.ClearStatus()
}

// UpdateFilename updates the output filename.
func (m *AppModel) UpdateFilename(filename string) {
	m.Filename = filename
	m.ClearStatus()
}

// ClearStatus clears status and error fields.
func (m *AppModel) ClearStatus() {
	m.Status = ""
	m.Error = nil
}

// Validate checks user inputs for correctness.
func (m *AppModel) Validate() error {
	if m.FontSizeText == "" {
		return errors.New("font size cannot be empty")
	}
	size, err := strconv.Atoi(m.FontSizeText)
	if err != nil {
		return fmt.Errorf("invalid font size format: %v", err)
	}
	if size <= 0 {
		return errors.New("font size must be a positive integer")
	}
	if m.Filename == "" {
		return errors.New("output filename cannot be empty")
	}
	return nil
}

// Generate validates inputs and orchestrates generating the image file.
func (m *AppModel) Generate() error {
	m.ClearStatus()
	if err := m.Validate(); err != nil {
		m.Error = err
		return err
	}

	size, _ := strconv.Atoi(m.FontSizeText)

	gr, err := util.OpenFont("goregular")
	if err != nil {
		m.Error = fmt.Errorf("error loading font: %v", err)
		return m.Error
	}

	grf := util.GetFontFace(float64(size), 180, gr)

	if err := phonenumber.DrawPhoneWithText(m.Result, m.Filename, grf); err != nil {
		m.Error = fmt.Errorf("error drawing image: %v", err)
		return m.Error
	}

	m.Status = fmt.Sprintf("Successfully generated image to %s", m.Filename)
	return nil
}
