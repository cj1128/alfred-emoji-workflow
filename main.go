package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fate-lovely/go-alfred"

	"os"
)

func main() {
	keyword := os.Args[1]
	for _, item := range search(keyword) {
		alfred.AddItem(alfred.Item{
			Title: item.desc,
			Arg:   codeToString(item.code),
			Icon: alfred.Icon{
				Path: fmt.Sprintf("imgs/%s.png", item.code),
			},
		})
	}
	json, _ := alfred.JSON()
	fmt.Print(json)
}

func codeToString(code string) string {
	n, _ := strconv.ParseInt(code[2:], 16, 32)
	return string(n)
}

func search(keyword string) []*Emoji {
	var result []*Emoji
	for _, item := range emojis {
		if strings.Index(item.keywords, keyword) != -1 {
			result = append(result, item)
		}
	}
	return result
}
