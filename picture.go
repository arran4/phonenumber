package phonenumber

import (
	"bytes"
	_ "embed"
	"fmt"
	wordwrap "github.com/arran4/golang-wordwrap"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
)

var fontFamily *canvas.FontFamily

var (
	//go:embed "phone.png"
	phoneImageBytes []byte
)

func DrawPhoneWithText(s string, fn string, fce font.Face) error {
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

	sw := wordwrap.NewSimpleWrapper(s, fce, wordwrap.HorizontalCenterLines)

	ls, pt, err := sw.TextToRect(phoneBounds, wordwrap.FitterIgnoreY{})
	height += float64(pt.Y)

	i := image.NewRGBA(image.Rect(0, 0, int(width), int(height)-phoneBounds.Max.Y))
	draw.Draw(i, i.Bounds(), image.White, image.Pt(0, 0), draw.Over)
	if err := sw.RenderLines(i, ls, i.Bounds().Min); err != nil {
		return fmt.Errorf("")
	}

	c := canvas.New(width, height)
	ctx := canvas.NewContext(c)

	ctx.SetFillColor(color.White)
	ctx.Fill()

	ctx.DrawImage((width-float64(i.Bounds().Dx()))/2, float64(phoneBounds.Dy()), i, canvas.Resolution(1))
	ctx.DrawImage((width-float64(phoneBounds.Dx()))/2, 0, phoneImage, canvas.Resolution(1))

	if err := renderers.Write(fn, c); err != nil {
		return fmt.Errorf("writing file %s: %w", fn, err)
	}
	return nil
}
