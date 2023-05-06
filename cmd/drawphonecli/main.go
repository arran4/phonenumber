package main

import (
	"flag"
	"fmt"
	"github.com/arran4/golang-wordwrap/util"
	"log"
	"os"
	"phonenumber"
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fn := fs.String("outfile", "out.png", "The output filename")
	if err := fs.Parse(os.Args); err != nil {
		fmt.Println("Error processing args", err)
	}
	if !fs.Parsed() {
		return
	}
	flag.Parse()
	s := phonenumber.Numbers("Hello how are you?", phonenumber.OpIgnoreSpace, phonenumber.OpDotPauses)
	log.Printf("'%s'", s)
	gr, err := util.OpenFont("goregular")
	if err != nil {
		fmt.Println("Error processing args", err)
		return
	}
	grf := util.GetFontFace(12, 180, gr)
	if err := phonenumber.DrawPhoneWithText(s, *fn, grf); err != nil {
		log.Printf("Error: %s", err)
		return
	}
	log.Printf("Wrote: %s", *fn)
}
