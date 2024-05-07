package main

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	// "github.com/alecthomas/repr"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"image/color"
	// "strings"
	"os"
	// "encoding/xml"
	"github.com/beevik/etree"	
	"log"
	"strings"
	// "reflect"
)
// var (
// 	// MVP regex for IPA form lexing
// 	IPAFormLexer = lexer.MustSimple([]lexer.SimpleRule{
// 		{`Token`, `kw|k|ɪ|ŋ|\s`},
// 	})
// 	parser = participle.MustBuild[Document](participle.Lexer(IPAFormLexer))
// )

type Document struct {
	Tokens []*Token `@@*`
}

type Token struct {
	Key string `@Token`
}

func (d *Document) PrintTokens() {
    for _, token := range d.Tokens {
        fmt.Println(token.Key)
    }
}

// func main() {
// 	ini, err := parser.Parse("", os.Stdin)
// 	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	
	// Create new canvas of dimension 100x100 mm
	c := canvas.New(800, 800)

	// Create a canvas context used to keep drawing state
	ctx := canvas.NewContext(c)
	var Transparent = color.RGBA{0x00, 0x00, 0x00, 0x00} // rgba(0, 0, 0, 0)	
	ctx.SetFillColor(Transparent)
	ctx.SetStrokeColor(canvas.Black)

	formsMap := make(map[string]string)
	
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("lang/system.svg"); err != nil {
		log.Fatalf("Failed to parse document: %v", err)
	}

	root := doc.SelectElement("svg")
	for _, path_element := range root.FindElements("//path") {
		label := path_element.SelectAttrValue("inkscape:label", "unknown")
		if strings.HasPrefix(label, "IPA:") {
			standard_IPA_prununciation  := strings.TrimPrefix(label, "IPA: ")
			SVG_path := path_element.SelectAttr("d").Value
			// TODO: remove the first 'move' command to make these paths relative
			formsMap[standard_IPA_prununciation] = SVG_path
		}
	}

	

	// for k, v := range formsMap {
	// 	fmt.Println("IPA: ", k)
	// 	fmt.Println("Path: ", v)
	// }

	// build regex for lexer based on user dictionary from the SVG
	var regexStrBuilder strings.Builder

	i := 0
	for k, _ := range formsMap {
		if i != 0 {
			regexStrBuilder.WriteString("|")
		}
		regexStrBuilder.WriteString(k)
		i++
	}

	regexStrBuilder.WriteString(`|\s`)

	fmt.Println("regex: ")
	fmt.Println(regexStrBuilder.String())

	// IPA lexer
	var (
		IPAFormLexer = lexer.MustSimple([]lexer.SimpleRule{
			{`Token`, regexStrBuilder.String()},
		})
		parser = participle.MustBuild[Document](participle.Lexer(IPAFormLexer)))


	// will be able to include longer strings when I finish the SVG dictionary
	// demo_string := "kwɪk kɪŋ" // ˈɡræmərz ɑr dɪˈfaɪnd æz tæɡd ɡoʊ ˈstrʌkʧərz. ˈpɑrtɪˌsɪpəl wɪl fɜrst lʊk fɔr tæɡz ɪn ðə fɔrm ˈpɑrsər ɪt wɪl ðɛn fɔl bæk tu ˈjuzɪŋ ði ɪnˈtaɪər tæɡ ˈbɑdi.`

	document, err := parser.Parse("", os.Stdin)
	// repr.Println(document, repr.Indent("  "), repr.OmitEmpty(true))
	// fmt.Println(document)
	document.PrintTokens()
	if err != nil {
		panic(err)
	}
	
	// point
	pos := canvas.Point{X: 10, Y: 100}

	// Simple sentence: "Welcome to a new way to write"
	message := [][]string{
		{"w", "ɛ", "l", "c", "ʌ", "m"}, 
		{"t"},
		{"ʌ"},
		{"n", "oo"},
		{"w", "eɪ"},
		{"t"},
		{"r", "aɪ", "t"},
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
