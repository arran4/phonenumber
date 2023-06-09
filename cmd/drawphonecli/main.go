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
