package phonenumber

import (
	"fmt"
	"github.com/arran4/golang-wordwrap/util"
)

// RunCLI is a subcommand `drawphonecli`
//
// Flags:
//
// 	outfile:  --outfile  (default: "out.png")         The output filename
// 	text:     --text     (default: "Hello how are you?") Text to write
// 	fontsize: --fontsize (default: 12)                Font size
func RunCLI(outfile string, text string, fontsize float64) error {
	s := Numbers(text, OpIgnoreSpace, OpDotPauses)
	fmt.Printf("'%s'\n", s)
	gr, err := util.OpenFont("goregular")
	if err != nil {
		return fmt.Errorf("error processing args: %w", err)
	}
	grf := util.GetFontFace(fontsize, 180, gr)
	if err := DrawPhoneWithText(s, outfile, grf); err != nil {
		return fmt.Errorf("error: %w", err)
	}
	fmt.Printf("Wrote: %s\n", outfile)
	return nil
}
