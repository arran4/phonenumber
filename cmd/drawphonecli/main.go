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
	"os"
	"phonenumber"
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fn := fs.String("outfile", "out.png", "The output filename")
	text := fs.String("text", "Hello how are you?", "Text to write")
	fontsize := fs.Float64("fontsize", 12, "Font size")
	if err := fs.Parse(os.Args); err != nil {
		fmt.Println("Error processing args", err)
	}
	if !fs.Parsed() {
		return
	}
	flag.Parse()
	s := phonenumber.Numbers(*text, phonenumber.OpIgnoreSpace, phonenumber.OpDotPauses)
	fmt.Printf("'%s'\n", s)
	gr, err := util.OpenFont("goregular")
	if err != nil {
		fmt.Println("Error processing args", err)
		return
	}
	grf := util.GetFontFace(*fontsize, 180, gr)
	if err := phonenumber.DrawPhoneWithText(s, *fn, grf); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Wrote: %s\n", *fn)
}
