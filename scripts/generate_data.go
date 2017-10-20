package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Emoji struct {
	Keywords []string `json:"keywords"`
	Char     string   `json:"char"`
	Imgid    string   `json:"imgid"`
}

const emojisJSONPath = "tmp/emojilib/emojis.json"
const outputPath = "emojis.go"

func main() {
	emojis, err := parseEmojisJSON()
	if err != nil {
		log.Fatal(err)
	}

	if err := generateEmojisGo(emojis); err != nil {
		log.Fatal(err)
	}
	fmt.Println("All done. 😉")
}

func parseEmojisJSON() (map[string]*Emoji, error) {
	var result map[string]*Emoji
	buf, err := ioutil.ReadFile(emojisJSONPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal json")
	}
	return result, nil
}

func generateEmojisGo(emojis map[string]*Emoji) error {
	out := bytes.Buffer{}

	out.WriteString(`package main

type Emoji struct {
  name string
  char string
  code string
  imgid string
  keywords string
}

var emojis = []*Emoji{
`)

	for name, e := range emojis {
		if e.Char == "" {
			continue
		}

		code := charToCode(e.Char)
		keywords := name + "," + strings.Join(e.Keywords, ",")

		line := fmt.Sprintf(
			"  {%s, %s, %s, %s, %s},\n",
			strconv.Quote(name),
			strconv.Quote(e.Char),
			strconv.Quote(code),
			strconv.Quote(e.Imgid),
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

func charToCode(char string) string {
	runes := []rune(char)
	result := make([]string, len(runes))
	for i, r := range runes {
		result[i] = fmt.Sprintf("%U", r)
	}
	return strings.Join(result, " ")
}
