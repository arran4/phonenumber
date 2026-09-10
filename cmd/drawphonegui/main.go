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

// Package main provides a GUI for generating Nokia phone number images.
package main

import (
	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
	"bitbucket.org/rj/goey/windows"
	"fmt"
	"os"
)

var (
	window *windows.Window
	model  *AppModel
)

func main() {
	model = NewAppModel()
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running GUI loop: %v\n", err)
		os.Exit(1)
	}
}

func createWindow() error {
	w, err := windows.NewWindow("Generate Nokia Phonenumbers", renderWindow())
	if err != nil {
		return err
	}
	w.SetScroll(false, true)
	window = w
	return nil
}

func updateWindow() {
	if window != nil {
		if err := window.SetChild(renderWindow()); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating window: %v\n", err)
		}
	}
}

func renderWindow() base.Widget {
	tabs := &goey.Tabs{
		Insets: goey.DefaultInsets(),
		Children: []goey.TabItem{
			renderTab(),
		},
		OnChange: func(_ int) {
			updateWindow()
		},
	}
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  tabs,
	}
}

func renderTab() goey.TabItem {
	var statusWidget base.Widget
	if model.Error != nil {
		statusWidget = &goey.Label{Text: fmt.Sprintf("Error: %v", model.Error)}
	} else if model.Status != "" {
		statusWidget = &goey.Label{Text: model.Status}
	} else {
		statusWidget = &goey.Label{Text: ""}
	}

	return goey.TabItem{
		Caption: "Configuration",
		Child: &goey.VBox{
			Children: []base.Widget{
				&goey.Label{Text: "Text:"},
				&goey.TextInput{
					Value:       model.Text,
					Placeholder: "Enter text to convert...",
					OnChange: func(v string) {
						model.UpdateText(v)
						updateWindow()
					},
					OnEnterKey: func(_ string) {
						updateWindow()
					},
				},
				&goey.Label{Text: "Result:"},
				&goey.TextInput{
					Value:       model.Result,
					Placeholder: "Result will appear here...",
					Disabled:    true,
				},
				&goey.Label{Text: "Font size:"},
				&goey.TextInput{
					Value:       model.FontSizeText,
					Placeholder: "e.g., 12",
					OnChange: func(v string) {
						model.UpdateFontSize(v)
						updateWindow()
					},
				},
				&goey.Label{Text: "Output filename:"},
				&goey.HBox{
					Children: []base.Widget{
						&goey.TextInput{
							Value:       model.Filename,
							Placeholder: "e.g., out.png",
							OnChange: func(v string) {
								model.UpdateFilename(v)
								updateWindow()
							},
						},
						&goey.Button{Text: "Browse...", OnClick: func() {
							browseFile()
						}},
					},
				},
				&goey.HBox{Children: []base.Widget{
					&goey.Button{Text: "Generate", Default: true, OnClick: func() {
						generate()
					}},
				}},
				statusWidget,
			},
		},
	}
}

func browseFile() {
	if window == nil {
		return
	}
	dlg := window.SaveFileDialog().
		WithTitle("Select Output File").
		WithFilename(model.Filename).
		AddFilter("PNG Images", "*.png").
		AddFilter("All Files", "*.*")

	filename, err := dlg.Show()
	if err != nil {
		// Ignore error since it might be cancellation or not supported
		return
	}
	if filename != "" {
		model.UpdateFilename(filename)
		updateWindow()
	}
}

func generate() {
	if err := model.Generate(); err != nil {
		if window != nil {
			_ = window.Message(fmt.Sprintf("Generation failed:\n%v", err)).
				WithTitle("Error").
				WithError().
				Show()
		}
	} else {
		if window != nil {
			_ = window.Message(fmt.Sprintf("Generation successful!\nSaved to: %s", model.Filename)).
				WithTitle("Success").
				WithInfo().
				Show()
		}
	}
	updateWindow()
}
