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

package phonenumber

import (
	"bytes"
	_ "embed"
	"fmt"
	wordwrap "github.com/arran4/golang-wordwrap"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"image/draw"
	"sync"
)

var (
	fontFamily     *canvas.FontFamily
	fontFamilyOnce sync.Once
	fontFamilyErr  error
)

var (
	//go:embed "phone.png"
	phoneImageBytes []byte
	phoneImage      canvas.Image
	phoneImageOnce  sync.Once
	phoneImageErr   error
)

func DrawPhoneWithText(s string, fn string, fce font.Face) error {
	fontFamilyOnce.Do(func() {
		fontFamily = canvas.NewFontFamily("times")
		if err := fontFamily.LoadFont(goregular.TTF, 0, canvas.FontRegular); err != nil {
			fontFamilyErr = err
		}
	})
	if fontFamilyErr != nil {
		return fmt.Errorf("loading font: %w", fontFamilyErr)
	}

	phoneImageOnce.Do(func() {
		var err error
		phoneImage, err = canvas.NewPNGImage(bytes.NewReader(phoneImageBytes))
		if err != nil {
			phoneImageErr = err
		}
	})
	if phoneImageErr != nil {
		return fmt.Errorf("loading phone image: %w", phoneImageErr)
	}
	phoneBounds := phoneImage.Bounds()
	width := float64(phoneBounds.Dx())
	height := float64(phoneBounds.Dy())

	sw := wordwrap.NewSimpleWrapper(s, fce, wordwrap.HorizontalCenterLines)

	ls, pt, _ := sw.TextToRect(phoneBounds, wordwrap.FitterIgnoreY{})
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
