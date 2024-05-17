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

// Load dictionary user-generated subforms and logograms
func load_script(path string) *Script {

	subforms := prefixtree.New()
	logograms := prefixtree.New()

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
			token := Token {
				Name: ipa_str,
					Path: canvas_path,
				}
			subforms.Add(ipa_str, token)
		}
		if logo_str := values.Get("logo"); logo_str != "" {
			canvas_path := SVG_path_to_canvas(path_element.SelectAttr("d").Value)
			token := Token {
				Name: logo_str,
					Path: canvas_path,
				}
			logograms.Add(logo_str, token)
		}
	}

	return &Script{
		SubForms: subforms,
		Logos: logograms,
	}
	
	// return tree, logograms
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
 Token can be one of the following 4 things:

   - alphabetical subform, representing a sequence of one or more IPA characters
   - logogram, representing a full word in the target language
   - phrase, representing multiple words in a sequence (not yet implemented)
   - whitespace, punctuation mark or other unknown symbol or character

   Naming: A subform is named the sequence of IPA characters it
   represents.  A logogram is named the word it represents in the
   target language (or in IPA at the discretion of the user). All other
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
 sequence of Tokens, for instance [r·aɪ·t·ɪŋ], representing the word
 "writing".

   Metaforms also contain extra information for debugging and
   rendering.
*/
type Metaform struct {
	Tokens []Token		// must have at least 1
	original_word string	// The string of characters represented by the Metaform pre-IPA conversion.
	// ipa_word		// IPA string
	contains_out_of_script_characters bool
	// Image - might store the rendered path here, not sure yet
	// Path
	height float64
	width float64
}

func (m Metaform) String() string {
	var tokens []string
	for _, token := range m.Tokens {
		tokens = append(tokens, token.Name)
	}
	return strings.Join(tokens, "·")
}


// A document is a sequence of Metaforms to be rendered.
type Document struct {
	Metaforms []Metaform
}

// returns a Document to be rendered
func Parse(input string, lcode string, script *Script) Document {
	// detach punctuation from words
	input = normalizePunctuation(input)

	ipa, err := LoadIPADict(lcode)
	if err != nil {
		log.Fatal(err)
	}

	words := strings.Fields(input)

	doc := Document{
		Metaforms: make([]Metaform, 0, len(words)*2),
	}

	// fmt.Println("logograms: ")
	// script.Logos.Output()
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
					original_word: matched_phrase_or_logo,
					// height:        10.0,
					// width:         15.0,
				}
				doc.Metaforms = append(doc.Metaforms, metaform)
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
			original_word: word,
			// contains_out_of_script_characters: false
			// height:        10.0,
			// width:         15.0,
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
		doc.Metaforms = append(doc.Metaforms, metaform)
		i++
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
	script := load_script("scripts/teen_script.svg")
	fmt.Println(script)
	// script.SubForms.Output()
	// script.Logos.Output()

	input_text := `when are you doing? how now brown cow? let's see how well we can do at testing logographs! This is not my forte, but I just want you to know about my system and what you can do with this`
	fmt.Println("input_text: ", input_text)
	// input_text := `Elephants, with their immense size and gracious movements, are a majestic sight in the wild.`
	// input_text := `this is just some writing`

	language_code := "en_US"

	document := Parse(input_text, language_code, script)
	fmt.Println(document)

	// pos := canvas.Point{X: 10, Y: 180}
	// yPos := pos.Y
	// for _, v := range document.Forms {
	// 	if v.Name == ` ` {
	// 		pos.Y = yPos
	// 		pos.X += 10
	// 		if pos.X >= 180 {
	// 			pos.X = 20
	// 			pos.Y -= 20
	// 			yPos = pos.Y
	// 		}
	// 	} else {

	// 	}
	// 	formPath, err := canvas.ParseSVGPath(v.Path)
	// 	if err == nil {
	// 		ctx.DrawPath(pos.X, pos.Y, formPath)
	// 		pos.X += formPath.Pos().X
	// 		pos.Y += formPath.Pos().Y
	// 	}
	// }

	// // Rasterize the canvas and write to a PNG file with 3.2 dots-per-mm (320x320 px)
	// if err := renderers.Write("rendered_text.png", c, canvas.DPMM(4)); err != nil {
	// 	panic(err)
	// }
}
