package main

import (
	"github.com/tdewolff/canvas"
	// "github.com/tdewolff/canvas/renderers/htmlcanvas"
	// "github.com/tdewolff/canvas/renderers"
	// "github.com/alecthomas/repr"
	// "github.com/alecthomas/participle/v2"
	// "github.com/alecthomas/participle/v2/lexer"
	"github.com/beevik/prefixtree"
	"fmt"
	// "io/ioutil"
	"encoding/json"
	"image/color"
	"github.com/beevik/etree"
	"log"
	"strconv"
	"strings"
	"net/url"
	"html"
	"errors"
	"syscall/js"
	"embed"
	// "unicode/utf8"
)


func renderWrapper() js.Func {
	renderFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 3 {
			return "Invalid no of arguments passed"
		}
		renderDoc := js.Global().Get("document")
		if !renderDoc.Truthy() {
			return "Unable to get document object"
		}
		renderOutputTextArea := renderDoc.Call("getElementById", "render_text_output")
		if !renderOutputTextArea.Truthy() {
			return "Unable to get output text area"
		}
		custom_script := args[0].Bool()
		script := args[1].String()

		fmt.Println("script: ", script)
		fmt.Println("custom_script (bool): ", custom_script)
		lcode := args[2].String()
		// fmt.Println("len(input)", len(inputSVG))
		status := render(custom_script, script, lcode)
		// if err != nil {
		//	fmt.Printf("unable to render using script %s\n", err)
		//	return err.Error()
		// }
		renderOutputTextArea.Set("value", status)
		return nil
		// return image
	})
	return renderFunc
}

// script is the SVG itself if the user uploaded a custom script,
// otherwise it's the name of one of the embedded scripts
func render(custom_script bool, script string, lcode string) string {
	// user supplied handwriting system definition
	s := load_script(custom_script, script)
	input_text := `How are you doing?

Let's see how well we can do at testing logographs! This is not my forte, but I just want you to know about my system and what you can do with this`
	fmt.Println("input_text: ", input_text)
	// language_code := "en_US"
	document := Parse(input_text, lcode, s)
	fmt.Println(draw_image(document))
	fmt.Println("document: ", document)
	return "success"
}


func draw_image(d Document) string {
	// Grab the canvas from the DOM
	cvs := js.Global().Get("document").Call("getElementById", "output_canvas")
	c := htmlcanvas.New(cvs, 200, 200, 4.0)
	// Create a canvas context used to keep drawing state
	ctx := canvas.NewContext(c)
	var Transparent = color.RGBA{0x00, 0x00, 0x00, 0x00} // Reba(0, 0, 0, 0)
	ctx.SetFillColor(Transparent)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(0.265)

	// space_between_metaforms := 10
	
	pos := canvas.Point{X: 10, Y: 180}
	yPos := pos.Y

	
	for _, m := range d.Metaforms {
		if m.Name == ` ` {
			pos.Y = yPos
			pos.X += 10
			if pos.X >= 180 {
				pos.X = 20
				pos.Y -= 20
				yPos = pos.Y
			}
		} else {
			for _, t := range m.Tokens {
				if t.Path != "" {
					formPath, err := canvas.ParseSVGPath(t.Path)
					if err == nil {
						ctx.DrawPath(pos.X, pos.Y, formPath)
						pos.X += formPath.Pos().X
						pos.Y += formPath.Pos().Y
					}
				}
			}
		}
	}
	return "success!!"
	// Create a triangle path from an SVG path and draw it to the canvas
	// triangle, err := canvas.ParseSVGPath("L60 0L30 60z")
	// if err != nil {
	//	panic(err)
	// }
	// ctx.SetFillColor(canvas.Mediumseagreen)
	// ctx.DrawPath(20, 20, triangle)
	// ctx.DrawPath(2, 2, "m 2 2 l 2 4")
}

func main() {

	fmt.Println("Go web assembly")
	js.Global().Set("renderSVG", renderWrapper())
	<-make(chan struct{})

	// user supplied handwriting system definition
	// script := load_script("scripts/teen_script.svg")
	// fmt.Println(script)
	// script.SubForms.Output()
	// script.Logos.Output()

	// input_text := `How are you doing? Let's see how well we can do at testing logographs! This is not my forte, but I just want you to know about my system and what you can do with this`
	// fmt.Println("input_text: ", input_text)
	// input_text := `Elephants, with their immense size and gracious movements, are a majestic sight in the wild.`
	// input_text := `this is just some writing`

	// language_code := "en_US"

	// document := Parse(input_text, language_code, script)
	// fmt.Println("document: ", document)

}
