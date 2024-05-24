package main

import (
	"log"
	"io"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"os"
)

func fileToString(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func main() {
	l := log.Default()
	options := fileToString("options.json")
	c := canvas.New(500, 500)
	// c := canvas.New
	ctx := canvas.NewContext(c)
	Render(ctx, options, l)

	if err := renderers.Write("rendered_text.png", c, canvas.DPMM(3.2)); err != nil {
		panic(err)
	}
	
	// f, err := os.Create("rendered_text.png")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// // Encode img to f in PNG format
	// err = png.Encode(f, img)
	// if err != nil {
	// 	panic(err)
	// }

}
