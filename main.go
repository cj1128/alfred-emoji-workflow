package main

import (
	"fmt"
	"sort"
	"strings"

	"os"

	"github.com/cj1128/gofred"
)

func main() {
	searchText := strings.TrimSpace(os.Args[1])

	if searchText == "" {
		searchText = "face"
	}

	terms := strings.Fields(searchText)

	result := search(terms)

	sort.SliceStable(result, func(i, j int) bool {
		return getScore(terms, result[i]) > getScore(terms, result[j])
	})

	for _, emoji := range result {
		iconPath := fmt.Sprintf("%s", emoji.img)

		gofred.AddItem(&gofred.Item{
			Title:    emoji.name,
			Subtitle: fmt.Sprintf(`Copy "%s" to clipboard`, emoji.char),
			Arg:      emoji.char,
			Valid:    true,
			Mods: gofred.Mods{
				gofred.CtrlKey: &gofred.Mod{
					Valid:    true,
					Arg:      ":" + emoji.name + ":",
					Subtitle: fmt.Sprintf(`Copy ":%s:" to clipboard`, emoji.name),
				},
				gofred.CmdKey: &gofred.Mod{
					Valid:    true,
					Arg:      emoji.img,
					Subtitle: "Open emoji image",
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

func search(terms []string) []*Emoji {
	var result []*Emoji

outer:
	for _, emoji := range emojis {
		for _, term := range terms {
			if strings.Index(emoji.keywords, term) == -1 {
				continue outer
			}
		}

		result = append(result, emoji)
	}

	return result
}

func getScoreForTerm(term string, emoji *Emoji) int {
	if term == emoji.name {
		return 4
	}

	i := strings.Index(emoji.name, term)

	if i == 0 {
		return 2
	}

	if i > 0 {
		return 1
	}

	return 0
}

func getScore(terms []string, emoji *Emoji) int {
	result := 0

	for _, term := range terms {
		result += getScoreForTerm(term, emoji)
	}

	return result
}
