package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
)

type Emoji struct {
	Name     string   `json:"name"`
	Char     string   `json:"char"`
	Keywords []string `json:"keywords"`

	Img string `json:"img"` // image path
}

type parsedResult struct {
	missingImageEmojis []*Emoji
	allEmojis          []*Emoji
}

const emojisJSONPath = "node_modules/emojilib/emojis.json"
const outputPath = "emojis.go"

func main() {
	result, err := parseEmojisJSON()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total emojis in emojilib: %d\n", len(result.allEmojis)+len(result.missingImageEmojis))
	fmt.Printf("Emojis with images: %d\n", len(result.allEmojis)-len(result.missingImageEmojis))
	fmt.Printf("Emojis missing images: %d\n", len(result.missingImageEmojis))

	// genMissingImageEmojisFile(result)

	if err := generateEmojisGo(result.allEmojis); err != nil {
		log.Fatal(err)
	}
	fmt.Println("All done. 😉")
}

func genMissingImageEmojisFile(result *parsedResult) {
	tmp, _ := os.Create("tmp.txt")
	w := tabwriter.NewWriter(tmp, 0, 1, 1, ' ', 0)
	defer tmp.Close()
	defer w.Flush()

	for _, e := range result.missingImageEmojis {
		fmt.Printf("  %s: %s\n", e.Name, e.Char)

		codePoints := charToCodePoints(e.Char)
		fmt.Fprintf(w, "%s\t;\t%s\t;\t%s\t\n", strings.Join(codePoints, " "), e.Name, e.Char)
	}
}

func parseEmojisJSON() (*parsedResult, error) {
	data := make(map[string]*Emoji)
	buf, err := ioutil.ReadFile(emojisJSONPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal emojis.json")
	}

	result := &parsedResult{}

	for name, value := range data {
		value.Name = name

		img := getEmojiImage(value.Char)

		if img == "" {
			result.missingImageEmojis = append(result.missingImageEmojis, value)
			value.Img = "placeholder.png"
		} else {
			value.Img = "imgs/" + img
		}
		result.allEmojis = append(result.allEmojis, value)
	}

	return result, nil
}

func generateEmojisGo(emojis []*Emoji) error {
	out := bytes.Buffer{}

	out.WriteString(`package main

type Emoji struct {
  name string
  char string
  img string
  keywords string
}

var emojis = []*Emoji{
`)

	for _, emoji := range emojis {
		keywords := emoji.Name + "," + strings.Join(emoji.Keywords, ",")

		line := fmt.Sprintf(
			"  {%s, %s, %s, %s},\n",
			strconv.Quote(emoji.Name),
			strconv.Quote(emoji.Char),
			strconv.Quote(emoji.Img),
			strconv.Quote(keywords),
		)
		out.WriteString(line)
	}
	out.WriteString("}")
	if err := writeData(out.Bytes()); err != nil {
		return err
	}
	return nil
}

func writeData(buf []byte) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "could not create data.go")
	}
	defer f.Close()
	_, err = f.Write(buf)
	if err != nil {
		return errors.Wrap(err, "could not write to file")
	}
	return nil
}

func charToCodePoints(char string) []string {
	runes := []rune(char)
	result := make([]string, len(runes))

	for i, r := range runes {
		result[i] = fmt.Sprintf("%x", r)
	}

	return result
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getEmojiImage(emoji string) string {
	codePoints := charToCodePoints(emoji)
	code := strings.Join(codePoints, "-")
	imgFilename := code + ".png"

	if fileExists(path.Join("node_modules", "emojiimages", "imgs", imgFilename)) {
		return imgFilename
	}

	// try add `fe0f` in the trailing
	imgFilename = code + "-fe0f.png"
	if fileExists(path.Join("node_modules", "emojiimages", "imgs", imgFilename)) {
		return imgFilename
	}

	// try add 'feof' in the middle
	if len(codePoints) == 2 {
		imgFilename = codePoints[0] + "-fe0f-" + codePoints[1] + ".png"
		if fileExists(path.Join("node_modules", "emojiimages", "imgs", imgFilename)) {
			return imgFilename
		}
	}

	return ""
}
