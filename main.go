package main

import (
	"fmt"
	"log"
	"syscall/js"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/htmlcanvas"
)

func renderWrapper() js.Func {
	renderFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		go func() {
			if len(args) != 1 {
				fmt.Println("Invalid no of arguments passed")
			}
			doc := js.Global().Get("document")
			if !doc.Truthy() {
				fmt.Println("Unable to get document object")
			}

			// Create a new logger that writes to the textarea
			logTextArea := doc.Call("getElementById", "log")
			if !logTextArea.Truthy() {
				fmt.Println("Unable to get log text area")
			}
			writer := jsWriter{logTextArea: logTextArea}
			logger := log.New(writer, "", 0)
			cvs := doc.Call("getElementById", "canvas")
			height := cvs.Get("height").Float()
			c := htmlcanvas.New(cvs, height, height, 1.0)
			ctx := canvas.NewContext(c)

			options := args[0].String()
			Render(ctx, options, logger)
		}()
		return nil
	})
	return renderFunc
}

// Define a custom writer
type jsWriter struct {
	logTextArea js.Value
}

// Implement the Write method for the io.Writer interface
func (jw jsWriter) Write(p []byte) (n int, err error) {
	// Append the text to the textarea, converting bytes to a JavaScript string
	currentText := jw.logTextArea.Get("value").String()
	newText := currentText + string(p)
	jw.logTextArea.Set("value", newText)

	// Return the number of bytes written and no error
	return len(p), nil
}

func main() {

	fmt.Println("Go web assembly")
	js.Global().Set("Render", renderWrapper())
	// js.Global().Set("Render", js.NewCallback)
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
