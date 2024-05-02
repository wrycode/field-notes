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
	c := canvas.New(200, 200)

	// Create a canvas context used to keep drawing state
	ctx := canvas.NewContext(c)
	var Transparent = color.RGBA{0x00, 0x00, 0x00, 0x00} // rgba(0, 0, 0, 0)	
	ctx.SetFillColor(Transparent)
	ctx.SetStrokeColor(canvas.Black)
	
	data, err := ioutil.ReadFile("forms.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// point
	pos := canvas.Point{X: 10, Y: 100}

	// Unmarshal the JSON data into a map
	var formsMap map[string]string
	err = json.Unmarshal(data, &formsMap)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Simple sentence: "Welcome to a new way to write"
	message := [][]string{
		{"w", "e", "l", "c", "u", "m"}, 
		{"t"},
		{"u"},
		{"n", "oo"},
		{"w", "ae"},
		{"t"},
		{"r", "wr_i_te", "t"},
	}

	// Render each word
	for _, word := range message {

		for _, form := range word {
			formPath, err := canvas.ParseSVGPath(formsMap[form])
			if err != nil {
				panic(err)
			}
			ctx.DrawPath(pos.X, pos.Y, formPath)
			pos.X += formPath.Pos().X
			pos.Y += formPath.Pos().Y
		}
		pos.X += 10
	}

	// Rasterize the canvas and write to a PNG file with 3.2 dots-per-mm (320x320 px)
	if err := renderers.Write("rendered_text.png", c, canvas.DPMM(3.2)); err != nil {
		panic(err)
	}
}
