package main

import (
	"fmt"
	"log"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"os"
	"io"
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
	c := canvas.New(500,500)
	l := log.Default()
	l.SetFlags(0)		// remove timestamp
	options := fileToString("options.json")
	Render(options, c, l)
	if err := renderers.Write("rendered_text.png", c, canvas.DefaultResolution); err != nil {
		fmt.Println("tes")
		panic(err)
	}

}
