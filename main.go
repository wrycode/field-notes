package main

import (
	"github.com/tdewolff/canvas"
	// "github.com/tdewolff/canvas/renderers"
	// "github.com/alecthomas/repr"
	// "github.com/alecthomas/participle/v2"
	// "github.com/alecthomas/participle/v2/lexer"
	"github.com/beevik/prefixtree"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"image/color"
	// "os"
	"github.com/beevik/etree"
	"log"
	"strconv"
	"strings"
	"net/url"
	"html"
)

// Load IPA dictionary from ipa_dicts/ dir as a map
// (from MIT-licensed https://github.com/open-dict-data/ipa-dict project)
func LoadIPADict(lang string) (map[string]string, error) {
	type IPAJson map[string][]map[string]string
	var jsonDict IPAJson

	file := fmt.Sprintf("./ipa_dicts/%s.json", lang)

	jsonFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonFile, &jsonDict)
	if err != nil {
		return nil, err
	}
	return jsonDict[lang][0], nil
}

// Convert SVG path description from Inkscape into path description for canvas library
func SVG_path_to_canvas(svgpath string) string {
	split := strings.Split(svgpath, " ")
	// remove the first 'move' command to make these paths relative
	split = split[2:]
	switch split[0] {
	case "m", "l", "h", "v", "c", "s", "q", "t", "a", "z":		// command is already present, don't need to do anything
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

// Load dictionary user-generated subforms and logograms
func load_system(path string) (*prefixtree.Tree, map[string]Form) {

	tree := prefixtree.New()
	logograms := make(map[string]Form)

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		log.Fatalf("Failed to parse document: %v", err)
	}

	root := doc.SelectElement("svg")
	for _, path_element := range root.FindElements("//path") {
		label := path_element.SelectAttrValue("inkscape:label", "unknown")

		values, err := url.ParseQuery(html.UnescapeString(label))
		if err != nil {
			panic(err)
		}
		if ipa_str := values.Get("IPA"); ipa_str != "" {
			canvas_path := SVG_path_to_canvas(path_element.SelectAttr("d").Value)
			form := Form {
				Name: ipa_str,
				Path: canvas_path,
				}
			tree.Add(ipa_str, form)
		}
		if logo := values.Get("logo"); logo != "" {
			canvas_path := SVG_path_to_canvas(path_element.SelectAttr("d").Value)
			form := Form {
				Name: logo,
				Path: canvas_path,
				}
			logograms[logo] = form
		}
	}
	return tree, logograms
}

// Return a string, placing a space before and after the punctuation
func normalizePunctuation(input string) string {
	punctuations := ".,;!?"
	for _, punc := range punctuations {
		// add spaces around punctuation only if they don't exist yet
		input = strings.ReplaceAll(input, string(punc), " "+string(punc)+" ")
	}

	// Remove extra spaces
	input = strings.Join(strings.Fields(input), " ")
	return input
}

// a Form is a 'unit' of handwritten lines - it can be an alphabetical
// 'subform', a punctuation mark, or a whole word (logogram). The path
// is a Canvas path. The pen can be picked up from the paper by
// issueing a MoveTo command (equivalent to SVG move). Some forms,
// like newlines, have no path value and are handled by the rendering
// code.
type Form struct {
	Name  string
	Path  string
}

func (f Form) String() string {
	return fmt.Sprintf("%v Path: %v", f.Name, f.Path)
}

// A document is a sequence of forms
type Document struct {
	Forms []Form
}

// func (d Document) String() string {
//     var forms []string
//     for _, form := range d.Forms {
//         forms = append(forms, form.String())
//     }
//     return strings.Join(forms, "\n")
// }

func main() {

	// Create new canvas of dimension 100x100 mm
	c := canvas.New(200, 200)

	// Create a canvas context used to keep drawing state
	ctx := canvas.NewContext(c)
	var Transparent = color.RGBA{0x00, 0x00, 0x00, 0x00} // Reba(0, 0, 0, 0)
	ctx.SetFillColor(Transparent)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(0.265)

	// user supplied handwriting system definition
	subformsTree, logo_map := load_system("lang/system.svg")

	fmt.Println(logo_map)
	fmt.Println(subformsTree)

	input_text := "question wrecks kick attack attach a catch net acclimation, holy day"
	lang := "en_US"

	// detach punctuation from words so it doesn't interfere with IPA or form lookup
	input_text = normalizePunctuation(input_text)

	// IPA dictionary
	ipa, err := LoadIPADict(lang)
	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(input_text)
	for i, word := range words {
		if replacement, exists := ipa[word]; exists {
			// when there are several possible
			// pronunciations, right now we just select
			// the first option
			first_option := strings.SplitN(replacement, ",", 2)[0]

			// strip forward slashes and accent characters
			// we're not using right now
			first_option = strings.ReplaceAll(first_option, "/", "")
			first_option = strings.ReplaceAll(first_option, "ˈ", "")
			first_option = strings.ReplaceAll(first_option, "ˌ", "")
			words[i] = first_option
		}
	}
	demo_ipa_string := strings.Join(words, " ")
	fmt.Print("demo_ipa_string: ")
	fmt.Println(demo_ipa_string)

	// point
	// pos := canvas.Point{X: 10, Y: 180}
	// yPos := pos.Y
	// for _, token := range document.Tokens {
	//	formPath, err := canvas.ParseSVGPath(formsMap[token.Key])
	//	if err == nil {
	//		ctx.DrawPath(pos.X, pos.Y, formPath)
	//		pos.X += formPath.Pos().X
	//		pos.Y += formPath.Pos().Y
	//	}
	//	if token.Key == ` ` {
	//		pos.Y = yPos
	//		pos.X += 10
	//		if pos.X >= 180 {
	//			pos.X = 20
	//			pos.Y -= 20
	//			yPos = pos.Y
	//		}
	//	}
	// }

	// // Rasterize the canvas and write to a PNG file with 3.2 dots-per-mm (320x320 px)
	// if err := renderers.Write("rendered_text.png", c, canvas.DPMM(4)); err != nil {
	//	panic(err)
	// }
}
