package main

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"image/color"
)

func main() {
	// Create new canvas of dimension 100x100 mm
	c := canvas.New(500, 500)

	// Create a canvas context used to keep drawing state
	ctx := canvas.NewContext(c)

	var Transparent = color.RGBA{0x00, 0x00, 0x00, 0x00} // rgba(0, 0, 0, 0)
	
	data, err := ioutil.ReadFile("forms.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	// Unmarshal the JSON data into a map
	var formsMap map[string]string
	err = json.Unmarshal(data, &formsMap)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Create an SVG path from the JSON form dictionary and draw it to the canvas
	form, err := canvas.ParseSVGPath(formsMap["l"])
	if err != nil {
		panic(err)
	}
	ctx.SetFillColor(Transparent)
	ctx.SetStrokeColor(canvas.Black)
	ctx.DrawPath(150, 150, form)

	form2, err := canvas.ParseSVGPath(formsMap["g"])
	if err != nil {
		panic(err)
	}
	ctx.SetFillColor(Transparent)
	ctx.SetStrokeColor(canvas.Black)
	ctx.DrawPath(176.30212, 150.23228, form2)
// 26.30212,0.23228
	// Rasterize the canvas and write to a PNG file with 3.2 dots-per-mm (320x320 px)
	if err := renderers.Write("rendered_text.png", c, canvas.DPMM(3.2)); err != nil {
		panic(err)
	}
}
