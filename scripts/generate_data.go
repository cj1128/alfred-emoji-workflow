// should run this script in parent dir with `make generate-data`
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v2"

	"github.com/pkg/errors"

	"github.com/PuerkitoBio/goquery"
)

const emojiListURL = "http://www.unicode.org/emoji/charts/emoji-list.html"

type ProgressReader struct {
	reader io.Reader
	bar    *pb.ProgressBar
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 {
		pr.bar.Add(n)
	}
	if err == io.EOF {
		pr.bar.Finish()
	}
	return n, err
}

func main() {
	fmt.Printf("Downloading %s...\n", emojiListURL)
	buf, err := downloadURL(emojiListURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsing...")
	if err := parseEmoji(buf); err != nil {
		log.Fatal(err)
	}
	fmt.Println("All done 😉 ")
}

func downloadURL(url string) ([]byte, error) {
	total, err := getContentLenght(url)
	if err != nil {
		return nil, errors.Wrap(err, "could not get content length")
	}

	bar := pb.Simple.Start(total)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	pr := &ProgressReader{resp.Body, bar}
	buf, err := ioutil.ReadAll(pr)
	if err != nil {
		return nil, errors.Wrap(err, "could not read buffer")
	}
	return buf, nil
}

func parseEmoji(buf []byte) error {
	out := bytes.Buffer{}

	out.WriteString(`package main

type Emoji struct {
  code string
  desc string
  keywords string
}

var emojis = []*Emoji{
`)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return errors.Wrap(err, "could not create new goquery document")
	}

	trs := doc.Find("tr")
	bar := pb.Simple.Start(trs.Length())

	trs.Each(func(i int, s *goquery.Selection) {
		defer bar.Increment()
		code := s.Find(".code").Text()
		if len(code) == 0 {
			return
		}

		// some emoji codes consist of two unicode points
		// e.g. U+1F476 U+1F3FE
		// we don't process emojis like this
		if strings.Index(code, " ") != -1 {
			return
		}

		names := s.Find("td.name")
		desc := names.First().Text()
		keywords := strings.Split(names.Last().Text(), " | ")
		imgBase64, _ := s.Find(".andr img").Attr("src")
		// remove 'data:image/png;base64,' prefix
		imgBase64 = imgBase64[22:]
		if err := saveImg(code, imgBase64); err != nil {
			log.Fatal(err)
		}

		line := fmt.Sprintf(
			"  {%s, %s, %s},\n",
			strconv.Quote(code),
			strconv.Quote(desc),
			strconv.Quote(strings.Join(append(keywords, desc), ", ")),
		)
		out.WriteString(line)
	})
	out.WriteString("}")
	if err := writeData(out.Bytes()); err != nil {
		return errors.Wrap(err, "could not write")
	}
	return nil
}

func saveImg(code, imgBase64 string) error {
	buf, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		return errors.Wrap(err, "could not decode base64")
	}
	name := fmt.Sprintf("workflow/imgs/%s.png", code)
	if err := ioutil.WriteFile(name, buf, 0644); err != nil {
		return errors.Wrap(err, "could not write file")
	}
	return nil
}

func writeData(buf []byte) error {
	f, err := os.Create("data.go")
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

func getContentLenght(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	str := resp.Header.Get("Content-Length")
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("could not convert Content-Length to integer: %s", str)
	}
	return n, nil
}
