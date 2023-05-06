package phonenumber

import (
	"bytes"
	_ "embed"
	"fmt"
	wordwrap "github.com/arran4/golang-wordwrap"
	"github.com/arran4/golang-wordwrap/util"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"image"
)

var fontFamily *canvas.FontFamily

var (
	//go:embed "phone.png"
	phoneImageBytes []byte
)

func DrawPhone(s string, fn string) error {
	fontFamily = canvas.NewFontFamily("times")
	if err := fontFamily.LoadLocalFont("NimbusRoman-Regular", canvas.FontRegular); err != nil {
		return fmt.Errorf("loading font: %w", err)
	}
	phoneImage, err := canvas.NewPNGImage(bytes.NewReader(phoneImageBytes))
	if err != nil {
		return fmt.Errorf("loading phone image: %w", err)
	}
	phoneBounds := phoneImage.Bounds()
	width := float64(phoneBounds.Dx())
	height := float64(phoneBounds.Dy())

	gr, err := util.OpenFont("goregular")
	if err != nil {
		return fmt.Errorf("font open error: %w", err)
	}
	grf := util.GetFontFace(16, 180, gr)

	sw := wordwrap.NewSimpleWrapper(s, grf)

	ls, pt, err := sw.TextToRect(phoneBounds, wordwrap.FitterIgnoreY{})
	height += float64(pt.Y)

	c := canvas.New(width, height)
	ctx := canvas.NewContext(c)

	ctx.DrawImage(0, height-float64(phoneBounds.Max.Y), phoneImage, canvas.DefaultResolution)

	i := image.NewRGBA(image.Rect(0, 0, int(width), int(height)-phoneBounds.Max.Y))
	if err := sw.RenderLines(i, ls, i.Bounds().Min); err != nil {
		return fmt.Errorf("")
	}

	if err := renderers.Write(fn, c); err != nil {
		return fmt.Errorf("writing file %s: %w", fn, err)
	}
	return nil
}
