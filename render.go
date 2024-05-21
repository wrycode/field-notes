package main

import (
	"net/http"
	"io"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
	"image"
	"github.com/beevik/prefixtree"
	"fmt"
	"encoding/json"
	"image/color"
	"github.com/beevik/etree"
	"log"
	"strconv"
	"strings"
	"net/url"
	"html"
	"errors"
	"embed"
	// "image/draw"
	// "image/png"
	// "os"

)

// Options for rendering
type Options struct {
	Script_debug bool `json:"Script_debug"`
	Custom_script_svg_value  string  `json:"Custom_script_svg_value"`
	Builtin_script_name string  `json:"Builtin_script_name"`
	Language_code string `json:"Language_code"`
	Image_width float64 `json:"Image_width"`
	Input_text string `json:"Input_text"`
}

func LoadIPADict(lang string) (map[string]string, error) {
	type IPAJson map[string][]map[string]string
	var jsonDict IPAJson

	resp, err := http.Get(fmt.Sprintf("https://storage.googleapis.com/ipa_dicts/%s.json", lang))
	if err != nil {
		fmt.Println("Error downloading IPA dictionary: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading IPA dictionar: " ,err)
	}

	err = json.Unmarshal(body, &jsonDict)
	if err != nil {
		fmt.Println("Error unmarshaling IPA dictionary: " ,err)
		return nil, err
	}
	return jsonDict[lang][0], err
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

/* A handwriting system definition

   SubForms are Tokens representing a sequence of one or more IPA
   characters. They are stored in a prefix tree so the parser can
   easily select the largest subform possible while scanning the IPA
   text.

   Logos are Tokens representing a full word (Logogram) or more than
   one word (phrase). They are also stored in a prefix tree.

*/

type Script struct {
	SubForms *prefixtree.Tree
	Logos *prefixtree.Tree
}

// Load dictionary of subforms and logograms
func load_script(script string) *Script {

	// TODO: need to validate the script to alert the user if an
	// IPA sequence (or word or phrase) is mapped to multiple
	// tokens. The other way around is fine, however: one token
	// can represent multiple IPA sequences (e.g. 'w' and 'hw' are
	// simplified to one symbol) or words or phrases.

	subforms := prefixtree.New()
	logograms := prefixtree.New()

	doc := etree.NewDocument()
	if err := doc.ReadFromString(script); err != nil {
		log.Fatalf("Failed to parse document: %v", err)
	}

	root := doc.SelectElement("svg")
	for _, path_element := range root.FindElements("//path") {
		label := path_element.SelectAttrValue("inkscape:label", "unknown")

		values, err := url.ParseQuery(html.UnescapeString(label))
		if err != nil {
			panic(err)
		}

		if ipa_field_val := values.Get("IPA"); ipa_field_val != "" {
			canvas_path := SVG_path_to_canvas(path_element.SelectAttr("d").Value)
			for _, ipa_sequence := range strings.Split(ipa_field_val, ",") {
				token := Token {
					Name: ipa_sequence,
						Path: canvas_path,
					}
				subforms.Add(ipa_sequence, token)
			}


		}

		if logo_field_val := values.Get("logo"); logo_field_val != "" {
			canvas_path := SVG_path_to_canvas(path_element.SelectAttr("d").Value)
			for _, logo_str := range strings.Split(logo_field_val, ",") {
				token := Token {
					Name: logo_str,
						Path: canvas_path,
					}
				logograms.Add(logo_str, token)
			}
		}
	}

	return &Script{
		SubForms: subforms,
		Logos: logograms,
	}
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
 specific script. The input text is parsed into a sequence of
 Tokens. A Token can be one of the following 3 things:

   - alphabetical subform, representing a sequence of one or more IPA characters
   - logogram, representing a full word or phrase in the target language
   - whitespace, punctuation mark or other unknown symbol or character

   Naming: A subform is named the sequence of IPA characters it
   represents.  A logogram is named the word or phrase it represents
   in the target language. All other Tokens represent a single
   character (whitespace or otherwise) and are named that character.

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
   connected together in a 'word'. A Metaform can be a single Token,
   like "," (comma), logograms, like "my" or "how are you", or a
   sequence of Tokens, for instance [r·aɪ·t·ɪŋ], representing the word
   "writing". ` ` (space) is the only token that is not converted to a
   Metaform, because all metaforms are assumed to be separated by a
   space. Spacing drawn between metaforms is therefore handled by the
   rendering code. Note that there IS a metaform for newlines which
   are handled by the rendering code, and all other whitespace chars (like
   tabs) are ignored.

   Metaforms also contain extra information for debugging and
   rendering.
*/
type Metaform struct {
	Name string // The string of characters represented by the Metaform
	Tokens []Token		// must have at least 1
	contains_out_of_script_characters bool //
	Img image.Image
	// Image - might store the rendered path here, not sure yet
	// Path - or here
	height float64
	width float64
}

func (m *Metaform) render(debug bool, width float64, height float64) *Metaform {
	// TODO: handle forms with no path
	// var path *canvas.Path
	// path, err := canvas.ParseSVGPath("")
	// if err != nil {
	//	fmt.Println("erro at line 243: ", err)
	// }

	// // Create a single path to work with
	// for _, t := range m.Tokens {
	//	if t.Path != "" {
	//		newpath, _ := canvas.ParseSVGPath(t.Path)
	//		path = path.Append(newpath)
	//	}
	// }
	// fmt.Println(m.Name)
	// fmt.Println("Path: ")
	// fmt.Println(path)


	// yPos := pos.Y

	c := canvas.New(width, height)
	ctx := canvas.NewContext(c)
	var Transparent = color.RGBA{0x00, 0x00, 0x00, 0x00} // Reba(0, 0, 0, 0)
	ctx.SetFillColor(Transparent)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(0.265)

	// Create a single path to work with
	pos := canvas.Point{X: 0, Y: 0}
	path, _ := canvas.ParseSVGPath("")

	for _, t := range m.Tokens {
		if t.Path != "" {
			newPath, err := canvas.ParseSVGPath(t.Path)
			if err != nil {
				fmt.Println("Error parsing path: for ",m.Name, t.Name, err)
			}
			newPath = newPath.Translate(pos.X, pos.Y)
			path = path.Join(newPath)
			pos.X = path.Pos().X
			pos.Y = path.Pos().Y
		}
	}

	ctx.DrawPath(0, 0, path)
	bounding_box := path.Bounds()
	ctx.SetStrokeColor(canvas.Red)
	ctx.SetStrokeWidth(0.1)

	ctx.DrawPath(0, 0, bounding_box.ToPath())

	// c.Clip(path.Bounds())
	c.Fit(1)

	// c.Fit()

	ctx.SetCoordRect(bounding_box, bounding_box.W, bounding_box.H)
	// m.height = 5.0
	m.Img = rasterizer.Draw(c, canvas.DefaultResolution, nil)


	return m
}


// func (m Metaform) String() string {
//	var tokens []string
//	for _, token := range m.Tokens {
//		tokens = append(tokens, token.Name)
//	}
//	return strings.Join(tokens, "·")
// }

// A document is a sequence of Metaforms to be rendered.
type Document struct {
	Metaforms []*Metaform
}

// returns a Document to be rendered
func Parse(input string, lcode string, script *Script) Document {
	// detach punctuation from words
	input = normalizePunctuation(input)

	input = strings.ToLower(input)

	ipa, err := LoadIPADict(lcode)
	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(input)

	doc := Document{
		Metaforms: make([]*Metaform, 0, len(words)*2),
	}

	// Loop through words, converting to logograms or IPA and appending to the document
	for i := 0; i < len(words); {
		word := words[i]

		// First check if the word is found in the logogram prefix tree
		val, err := script.Logos.FindValue(word)
		var matched_phrase_or_logo string // empty unless we find a matching sequence of one or more words in the script.Logos prefix tree

		if err == nil {
			logo := val.(Token)
			next_word_pos := i + 1 // Only used to update i later if we find a matching phrase or logo

			if logo.Name == word  { // exact match
				matched_phrase_or_logo = word
			}

			// Keep checking words until the sequence of
			// words isn't found in the prefix tree
			for j := i + 2; j < len(words); j++ {
				next_phrase := strings.Join(words[i:j], " ")
				// fmt.Println("next_phrase: ", next_phrase)
				next_val, err := script.Logos.FindValue(next_phrase)

				if err == nil {
					next_logo := next_val.(Token)
					// fmt.Println("1st")
					if next_logo.Name == next_phrase { // exact match
						logo = next_logo
						matched_phrase_or_logo = next_phrase
						next_word_pos = j
					}
				} else if errors.Is(err, prefixtree.ErrPrefixAmbiguous) {
					// fmt.Println("2nd")
					// continue searching for a phrase but don't save the current one
					// because we don't have a new match
				} else { // No match for the current word at j; time to save the logogram
					// fmt.Println("3rd")
					break
				}
			}

			if matched_phrase_or_logo != "" {
				// fmt.Println
				i = next_word_pos
				metaform := Metaform{
					Tokens:        []Token{logo},
					Name: matched_phrase_or_logo,
					// height:        10.0,
					// width:         15.0,
				}
				doc.Metaforms = append(doc.Metaforms, &metaform)
				continue
			}
		}

		// If we got this far in the loop, then we can assume
		// no phrase or logogram was found. Convert the word
		// to IPA:
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

		// Now convert the IPA characters into subform tokens

		// easier to index and loop through, we just
		// have to cast back into string when we
		// search the subforms prefix tree
		chars := []rune(words[i])

		end := len(chars)
		current_char := 0

		metaform := Metaform{
			Tokens:        []Token{},
			Name: word,
		}

		for current_char < end {
			seq_end := current_char + 1
			form_key := string(chars[current_char:seq_end]) // default to one character form
			val, err := script.SubForms.FindValue(form_key)

			if errors.Is(err, prefixtree.ErrPrefixNotFound) {
				// character is not defined in the script, so we'll just store it as a pathless Token and move on
				token := Token {
					Name: form_key,
					}
				metaform.Tokens = append(metaform.Tokens, token)
				metaform.contains_out_of_script_characters = true
			}

			for seq_end < end {
				seq_end+= 1
				new_form_key := string(chars[current_char:seq_end])
				new_val, err := script.SubForms.FindValue(new_form_key)

				if err == nil { // new, longer form found
					if new_val.(Token).Name == new_form_key  { // exact match
						form_key = new_form_key
						val = new_val
					}
				} else if errors.Is(err, prefixtree.ErrPrefixAmbiguous) {
				} else {
					seq_end -= 1 // need
					// to backtrack one character since no match was found
					break
				}
			}
			if val != nil {
				metaform.Tokens = append(metaform.Tokens, val.(Token))
			}
			current_char = seq_end
		}
		doc.Metaforms = append(doc.Metaforms, &metaform)
		i++
	}

	return doc
}

//go:embed scripts/*
var scripts embed.FS

// Renders the handwritten output to the provided canvas
func Render(options string, c *canvas.Canvas, log *log.Logger) {
	var o Options
	// Unmarshal the JSON string into opts
	err := json.Unmarshal([]byte(options), &o)

	if err != nil {
		log.Println(err)
	}

	// fmt.Println("printing o")
	// log.Println(o)

	// handwriting system definition
	var script *Script
	if o.Custom_script_svg_value != "" {
		script = load_script(o.Custom_script_svg_value)
	} else {
		file := fmt.Sprintf("scripts/%s", o.Builtin_script_name)
		script_file, err := scripts.ReadFile(file)
		if err != nil {
			log.Println("Error reading builtin script: ", err)
		}
		script = load_script(string(script_file))
	}

	if err != nil {
		log.Println("Error reading script: ", err)
	}

	// document to render
	d := Parse(o.Input_text, o.Language_code, script)

	// Create a canvas context used to keep drawing state
	ctx := canvas.NewContext(c)

	fmt.Println(ctx)
	for _, m := range d.Metaforms {
		m = m.render(false, 50, 50)
	}

	// fmt.Println(d.Metaforms)
	// fmt.Println(d.Metaforms[0])
	pos := canvas.Point{X: 10, Y: 450}


	for _, m := range d.Metaforms {
		if pos.X >= 450 {
			pos.X = 10
			pos.Y -= 90
		}
		// image := m.Img
		fmt.Println("Name: ", m.Name)
		ctx.DrawImage(pos.X, pos.Y, m.Img, 1)
		pos.X += float64(m.Img.Bounds().Dx() + 10)
		// pos.Y +=

		// pos.X += m.Img.Rect.Max.X
		// pos.Y += m.Img.Rect.Max.Y
		// fmt.Println(image)
	}

	// func (c *Context) DrawImage(x, y float64, img image.Image, resolution Resolution)

	// start stitching
	// upperLeftX := d.Metaforms[0].Img.Bounds().Max.X
	// upperLeftY := d.Metaforms[0].Img.Bounds().Max.Y
	// upperLeftX := 5
	// upperLeftY := 5

	// imgWidth := upperLeftX * len(d.Metaforms)
	// imgHeight := upperLeftY

	// // create new blank image with a size that depends on number of images
	// newImage := image.NewNRGBA(image.Rect(0, 0, 100, 100))

	// // Start drawing images from the slice to new blank image
	// for i, m := range d.Metaforms {
	//	if m.Img != nil {
	//	rect := image.Rect(upperLeftX*i, 0, upperLeftX*i+upperLeftX, imgHeight)
	//		draw.Draw(newImage, rect, m.Img, m.Img.Bounds().Min, draw.Src)

	//	}
	// }

	// // Create resulting image file on disk.
	// imgFile, err := os.Create("tiled.png")
	// if err != nil {
	//	panic(err)
	// }
	// defer imgFile.Close()

	// // Encode writes the Image m to w in PNG format.
	// err = png.Encode(imgFile, newImage)
	// if err != nil {
	//	panic(err)
	// }


}
