package main

import (
	"fmt"
	"log"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func main() {

	// Create new canvas
	c := canvas.New(500, 500)

	fmt.Println()
	l := log.Default()
	options := `{
        "Script_debug": false,
	"Custom_script_svg_value":  "",
        "Builtin_script_name": "teen_script.svg",
        "Language_code": "en_US",
        "Image_width": 500,
        "Input_text": "How are you doing? Let's see how well we can do at testing logographs! This is not my forte, but I just want you to know about my system and what you can do with this"
	}`
	fmt.Println("hello1")

	Render(options, c, l)

	// Rasterize the canvas and write to a PNG file with 3.2 dots-per-mm (320x320 px)
	if err := renderers.Write("rendered_text.png", c, canvas.DPMM(3.2)); err != nil {
		fmt.Println("tes")
		panic(err)
	}

}
