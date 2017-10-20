package main

import (
	"fmt"
	"strings"

	"github.com/fate-lovely/go-alfred"

	"os"
)

func main() {
	var result []*Emoji
	keyword := strings.TrimSpace(os.Args[1])
	if keyword == "" {
		keyword = "face"
	}
	result = search(keyword)
	for _, item := range result {
		iconPath := fmt.Sprintf("imgs/%s.png", item.imgid)
		alfred.AddItem(alfred.Item{
			Title:    item.name,
			Subtitle: fmt.Sprintf(`Copy "%s" to clipboard`, item.char),
			Arg:      item.char,
			Mods: alfred.Mods{
				"alt": alfred.Mod{
					Valid:    true,
					Arg:      ":" + item.name + ":",
					Subtitle: fmt.Sprintf(`Copy ":%s:" to clipboard`, item.name),
				},
			},
			Icon: alfred.Icon{
				Path: iconPath,
			},
			Text: alfred.Text{
				Copy:      item.char,
				Largetype: item.char,
			},
		})
	}
	json, _ := alfred.JSON()
	fmt.Print(json)
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
