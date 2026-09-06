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

// Package main provides a CLI for generating Nokia phone number images.
package main

import (
	"flag"
	"fmt"
	"github.com/arran4/golang-wordwrap/util"
	"io"
	"os"
	"phonenumber"
)

func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("drawphonecli", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fn := fs.String("outfile", "out.png", "The output filename")
	text := fs.String("text", "Hello how are you?", "Text to write")
	fontsize := fs.Float64("fontsize", 12, "Font size")

	if err := fs.Parse(args); err != nil {
		return err
	}

	s := phonenumber.Numbers(*text, phonenumber.OpIgnoreSpace, phonenumber.OpDotPauses)
	fmt.Fprintf(stdout, "'%s'\n", s)

	gr, err := util.OpenFont("goregular")
	if err != nil {
		return fmt.Errorf("loading font: %w", err)
	}
	grf := util.GetFontFace(*fontsize, 180, gr)

	if err := phonenumber.DrawPhoneWithText(s, *fn, grf); err != nil {
		return fmt.Errorf("drawing phone: %w", err)
	}
	fmt.Fprintf(stdout, "Wrote: %s\n", *fn)

	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
