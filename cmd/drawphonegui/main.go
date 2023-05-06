package main

import (
	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
	"bitbucket.org/rj/goey/windows"
	"fmt"
	"github.com/arran4/golang-wordwrap/util"
	"phonenumber"
	"strconv"
)

var (
	window   *windows.Window
	fn       = "out.png"
	text     = "Hello how are you?"
	result   = phonenumber.Numbers(text, phonenumber.OpIgnoreSpace, phonenumber.OpDotPauses)
	fontsize = "12"
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
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
	window.SetChild(renderWindow())
}

func renderWindow() base.Widget {
	tabs := &goey.Tabs{
		Insets: goey.DefaultInsets(),
		Children: []goey.TabItem{
			renderTab(),
		},
		OnChange: func(i int) {
			updateWindow()
		},
	}
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  tabs,
	}
}

func renderTab() goey.TabItem {
	return goey.TabItem{
		Caption: "Configuration / Authentication",
		Child: &goey.VBox{
			Children: []base.Widget{
				&goey.Label{Text: "Text:"},
				&goey.TextInput{
					Value:       text,
					Placeholder: "Hidden",
					OnChange: func(v string) {
						text = v
						result = phonenumber.Numbers(v, phonenumber.OpIgnoreSpace, phonenumber.OpDotPauses)
					},
					OnEnterKey: func(value string) {
						updateWindow()
					},
				},
				&goey.Label{Text: "Result: (Press enter above to refresh)"},
				&goey.TextInput{
					Value:       result,
					Placeholder: "Hidden",
					Disabled:    true,
				},
				&goey.Label{Text: "Font size:"},
				&goey.TextInput{
					Value:       fontsize,
					Placeholder: "Hidden",
					OnChange: func(v string) {
						fontsize = v
						updateWindow()
					},
				},
				&goey.Label{Text: "Output filename:"},
				&goey.TextInput{
					Value:       fn,
					Placeholder: "Hidden",
					OnChange: func(v string) {
						fn = v
						updateWindow()
					},
				},
				&goey.HBox{Children: []base.Widget{
					&goey.Button{Text: "Generate", Default: true, OnClick: func() {
						generate()
					}},
				}},
			},
		},
	}
}

func generate() {
	s := phonenumber.Numbers(text, phonenumber.OpIgnoreSpace, phonenumber.OpDotPauses)
	result = s
	gr, err := util.OpenFont("goregular")
	if err != nil {
		fmt.Println("Error processing args", err)
		return
	}
	atoi, _ := strconv.Atoi(fontsize)
	grf := util.GetFontFace(float64(atoi), 180, gr)
	if err := phonenumber.DrawPhoneWithText(s, fn, grf); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Wrote: %s\n", fn)
}
