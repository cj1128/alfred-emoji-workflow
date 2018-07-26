package main

import (
	"fmt"
	"sort"
	"strings"

	"os"

	"github.com/fate-lovely/gofred"
)

func main() {
	keyword := strings.TrimSpace(os.Args[1])

	if keyword == "" {
		keyword = "face"
	}

	result := search(keyword)

	sort.SliceStable(result, func(i, j int) bool {
		return getScore(keyword, result[i]) > getScore(keyword, result[j])
	})

	for _, emoji := range result {
		iconPath := fmt.Sprintf("imgs/%s", emoji.img)

		gofred.AddItem(&gofred.Item{
			Title:    emoji.name,
			Subtitle: fmt.Sprintf(`Copy "%s" to clipboard`, emoji.char),
			Arg:      emoji.char,
			Valid:    true,
			Mods: gofred.Mods{
				gofred.CmdKey: &gofred.Mod{
					Valid:    true,
					Arg:      ":" + emoji.name + ":",
					Subtitle: fmt.Sprintf(`Copy ":%s:" to clipboard`, emoji.name),
				},
			},
			Icon: &gofred.Icon{
				Path: iconPath,
			},
			Text: &gofred.Text{
				Copy:      emoji.char,
				Largetype: emoji.char,
			},
		})
	}

	json, _ := gofred.JSON()
	fmt.Print(json)
}

func search(keyword string) []*Emoji {
	var result []*Emoji
	for _, emoji := range emojis {
		if strings.Index(emoji.keywords, keyword) != -1 {
			result = append(result, emoji)
		}
	}
	return result
}

func getScore(keyword string, emoji *Emoji) int {
	if keyword == emoji.name {
		return 3
	}

	i := strings.Index(emoji.name, keyword)
	if i == 0 {
		return 2
	}
	if i > 0 {
		return 1
	}

	return 0
}
