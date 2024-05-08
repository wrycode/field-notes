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
	// "os"
	// "encoding/xml"
	"github.com/beevik/etree"
	"log"
	"strconv"
	"strings"
	// "reflect"
)
// var (
//	// MVP regex for IPA form lexing
//	IPAFormLexer = lexer.MustSimple([]lexer.SimpleRule{
//		{`Token`, `kw|k|ɪ|ŋ|\s`},
//	})
//	parser = participle.MustBuild[Document](participle.Lexer(IPAFormLexer))
// )

type Document struct {
	Tokens []*Token `@@*`
}

type Token struct {
	Key string `@Token`
}

func (d *Document) PrintTokens() {
	for _, token := range d.Tokens {
		fmt.Print(token.Key, " ")

	}
	fmt.Println()
}

// func main() {
//	ini, err := parser.Parse("", os.Stdin)
//	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
//	if err != nil {
//		panic(err)
//	}
// }

func SVG_path_to_canvas(svgpath string) string {

	split := strings.Split(svgpath, " ")
	// remove the first 'move' command to make these paths relative
	split = split[2:]
	switch split[0] {
	case "m", "l", "h", "v", "c", "s", "q", "t", "a", "z":
		// don't need to do anything
	default:
		split = append([]string{"l"}, split...) // need to add the lineto command
	}

	for i, str := range split {
		split[i] = reverseYAxis(str)
	}
	return strings.Join(split, " ")
}

func reverseYAxis(coord string) string {
	points := strings.Split(coord, ",")
	if len(points) != 2 {
		return coord
	}
	y, err := strconv.ParseFloat(points[1], 64)
	if err != nil {
		return coord
	}
	y = -y
	return fmt.Sprintf("%s,%f", points[0], y)
}


func main() {

	// Create new canvas of dimension 100x100 mm
	c := canvas.New(200, 200)

	// Create a canvas context used to keep drawing state
	ctx := canvas.NewContext(c)
	var Transparent = color.RGBA{0x00, 0x00, 0x00, 0x00} // rgba(0, 0, 0, 0)
	ctx.SetFillColor(Transparent)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(0.265)

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
			formsMap[standard_IPA_prununciation] = SVG_path_to_canvas(SVG_path)
		}
	}


	for k, v := range formsMap {
		fmt.Println("IPA: ", k)
		fmt.Println("Path: ", v)
	}

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

	regexStrBuilder.WriteString(`|\s|.`)

	// fmt.Println("regex: ")
	// fmt.Println(regexStrBuilder.String())

	// IPA lexer
	var (
		IPAFormLexer = lexer.MustSimple([]lexer.SimpleRule{
			{`Token`, regexStrBuilder.String()},
		})
		parser = participle.MustBuild[Document](participle.Lexer(IPAFormLexer)))

	demo_string := `heɪ, haʊz ɪt ˈɡoʊɪŋ? aɪ ʤʌst keɪm frʌm ə ˈkreɪzi deɪ æt wɜrk. ju woʊnt bɪˈliv wɑt ˈhæpənd. ˈaʊər bɑs ˈsʌdənli ˌdɪˈsaɪdɪd ðæt wi nid ə ˈtoʊtəl riˈvæmp fɔr ˈaʊər ˈprɑʤɛkt. naʊ ˈɪzənt ðæt ʤʌst ˈpiʧi? aɪ min, wiv bɪn ˈwɜrkɪŋ ɑn ˈɡɛtɪŋ ðoʊz dræfts dʌn fɔr wiks!`

	// document, err := parser.Parse("", os.Stdin)
	document, err := parser.ParseString("", demo_string)
	// repr.Println(document, repr.Indent("  "), repr.OmitEmpty(true))
	// fmt.Println(document)
	document.PrintTokens()
	if err != nil {
		panic(err)
	}

	// Create a triangle path from an SVG path and draw it to the canvas
	triangle, err := canvas.ParseSVGPath("L0.6 0L0.3 0.6z")
	if err != nil {
		panic(err)
	}
	ctx.SetFillColor(canvas.Mediumseagreen)
	ctx.DrawPath(30, 180, triangle)
	ctx.SetFillColor(Transparent)

	// point
	pos := canvas.Point{X: 10, Y: 180}
	yPos := pos.Y
	for _, token := range document.Tokens {
		formPath, err := canvas.ParseSVGPath(formsMap[token.Key])

		if err == nil {
			ctx.DrawPath(pos.X, pos.Y, formPath)
			pos.X += formPath.Pos().X
			pos.Y += formPath.Pos().Y
		}
		if token.Key == ` ` {
			pos.Y = yPos
			pos.X += 10
			if pos.X >= 180 {
				pos.X = 20
				pos.Y -= 20
				yPos = pos.Y
			}




		}


	}


	// // Simple sentence: "Welcome to a new way to write"
	// message := [][]string{
	//	{"w", "ɛ", "l", "c", "ʌ", "m"},
	//	{"t"},
	//	{"ʌ"},
	//	{"n", "oo"},
	//	{"w", "eɪ"},
	//	{"t"},
	//	{"r", "aɪ", "t"},
	// }

	// // Render each word
	// for _, word := range message {

	//	for _, form := range word {
	//		formPath, err := canvas.ParseSVGPath(formsMap[form])
	//		if err != nil {
	//			panic(err)
	//		}
	//		ctx.DrawPath(pos.X, pos.Y, formPath)
	//		pos.X += formPath.Pos().X
	//		pos.Y += formPath.Pos().Y
	//	}
	//	pos.X += 10
	// }

	// Rasterize the canvas and write to a PNG file with 3.2 dots-per-mm (320x320 px)
	if err := renderers.Write("rendered_text.png", c, canvas.DPMM(4)); err != nil {
		panic(err)
	}
}
