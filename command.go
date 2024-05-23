package main

import (
	// "fmt"
	"log"
	"io"
	// "image"
	// "image/color"
	"image/png"
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
	img := Render(options, l)

	f, err := os.Create("rendered_text.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Encode img to f in PNG format
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}

}
