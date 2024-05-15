package main

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
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
	"errors"
	// "unicode/utf8"
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
// TODO: May need to improve this function in the future! Have not tested extensively.
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
func load_script(path string) (*prefixtree.Tree, map[string]Form) {

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

/* a Token is an indivisible 'unit' of a handwritten document in a
 specific script. The input text is parsed into a sequence of Tokens. A
 Token can be one of the following 5 things:

   - alphabetical subform, representing a sequence of one or more IPA characters
   - logogram, representing a full word in the target language
   - phrase, representing multiple words in a sequence (not yet implemented)
   - whitespace
   - punctuation mark or other unknown symbol or character

   Naming: A subform is named the sequence of IPA characters it
   represents.  A logogram is named the word it represents in the
   target language or in IPA at the discretion of the user. All other
   Tokens represent a single character (whitespace or otherwise) and
   are named that character.

   If the Token is defined in the current Script, Path contains a
   Canvas path, otherwise it is empty.
*/

type Token struct {
	Name  string
	Path  string
}

func (t Token) String() string {
	return fmt.Sprintf("%v Path: %v", t.Name, t.Path)
}

/* A Metaform is a sequence of one or more Tokens that are drawn
 connected together. A Metaform can be a single Token, like ","
 (comma), " " (space), "my" (logogram for the word "my"), or a
 sequence of Tokens, for instance [r·a·ɪ·t·ɪ·ŋ], representing the word
 "writing".

   Metaforms also contain extra information for debugging and
   rendering.
*/
type Metaform struct {
	Forms []Form		// must have at least 1
	original_word string	// The string of characters represented by the Metaform pre-IPA conversion
	// Image - might store the rendered path here, not sure yet
	height float64
	width float64
}

// A document is a sequence of Metaforms to be rendered.
type Document struct {
	Metaforms []Metaform
}


// func (d Document) String() string {
//	var forms []string
//	for _, form := range d.Forms {
//		forms = append(forms, form.Name)
//	}
//	return strings.Join(forms, "·")
// }

// returns a Document to be rendered
func Parse(input string, lcode string, subforms *prefixtree.Tree, logos map[string]Form) Document {

	// detach punctuation from words so it doesn't interfere with IPA or form lookup
	input = normalizePunctuation(input)

	ipa, err := LoadIPADict(lcode)
	if err != nil {
		log.Fatal(err)
	}

	// convert string to IPA, skipping logograms which are not defined in IPA
	words := strings.Fields(input)
	for i, word := range words {
		if _, ok := logos[word]; !ok { // skip word if it is in logos map
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
	}

	ipa_string := strings.Join(words, " ")
	fmt.Print("ipa_string: ")
	fmt.Println(ipa_string)

	doc := Document{
		Forms: make([]Form, 0, len(words)*2),
	}

	// parse each word into subforms or 1 logogram and add to the document
	for _, word := range words {
		if val, ok := logos[word]; ok {
			doc.Forms = append(doc.Forms, val)
		} else {		// turn word into subforms

			// easier to index and loop through, we just
			// have to cast back into string when we
			// search the subforms prefix tree
			chars := []rune(word)

			end := len(chars)
			current_char := 0

			for current_char < end {
				seq_end := current_char + 1
				form_key := string(chars[current_char:seq_end]) // default to one character form
				val, _ := subforms.FindValue(form_key) // TODO check error?
				// fmt.Print("one char form search value and error: ")
				// fmt.Println(val, err)

				for seq_end < end {
					seq_end+= 1
					new_form_key := string(chars[current_char:seq_end])
					new_val, err := subforms.FindValue(new_form_key)
					// fmt.Print("new form search value and error: ")
					// fmt.Println(new_form_key, new_val, err)

					if err == nil { // new, longer form found
						// fmt.Println("1st")

						if new_val.(Form).Name == new_form_key  { // exact match
							form_key = new_form_key
							val = new_val
						}
					} else if errors.Is(err, prefixtree.ErrPrefixAmbiguous) {
						// fmt.Println("2nd")

					} else {
						seq_end -= 1 // need
						// to backtrack one character since no match was found
						// fmt.Println("3rd")
						break
					}
				}
				if val != nil {
					doc.Forms = append(doc.Forms, val.(Form))
				}
				current_char = seq_end
			}

		}

		// Add a space between words, which does not have a path value because we will handle rendering
		doc.Forms = append(doc.Forms, Form {
			Name: " ",
				Path: "dummy",
			})
	}

	return doc

}

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
	script_subforms, script_logograms := load_script("scripts/demotic.svg")

	// input_text := `let's see how well we can do at testing logographs! This is not my forte, but I just want you to know about my system and what you can do with this`
	// input_text := `Elephants, with their immense size and gracious movements, are a majestic sight in the wild.`
	input_text := `this is just some writing`

	language_code := "en_US"

	document := Parse(input_text, language_code, script_subforms, script_logograms)
	fmt.Println(document)


	pos := canvas.Point{X: 10, Y: 180}
	yPos := pos.Y
	for _, v := range document.Forms {
		if v.Name == ` ` {
			pos.Y = yPos
			pos.X += 10
			if pos.X >= 180 {
				pos.X = 20
				pos.Y -= 20
				yPos = pos.Y
			}
		} else {

		}
		formPath, err := canvas.ParseSVGPath(v.Path)
		if err == nil {
			ctx.DrawPath(pos.X, pos.Y, formPath)
			pos.X += formPath.Pos().X
			pos.Y += formPath.Pos().Y
		}
	}

	// Rasterize the canvas and write to a PNG file with 3.2 dots-per-mm (320x320 px)
	if err := renderers.Write("rendered_text.png", c, canvas.DPMM(4)); err != nil {
		panic(err)
	}
}
