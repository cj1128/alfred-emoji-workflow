# Alfred Emoji Workflow 🙈

[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://mit-license.org/2016)

Alfred workflow for quick searching emojis with keyword. After all, every one loves emoji 😎.

<p align="center">
  <img src="http://ww1.sinaimg.cn/large/9b85365dgy1ftn2byughhg20mg0fc7wj" />
</p>

## Installation 😆

click to download latest version of [Emoji.alfredworkflow](https://github.com/cj1128/alfred-emoji-workflow/releases/download/v1.2.0/Emoji.alfredworkflow).

## Usage 🌟

- press `Enter` to copy emoji, e.g. 😎
- press `Ctrl + Enter` to copy text representation, e.g. `:sunglasses:`
- press `Cmd + Enter` to open emoji image

## Development

This workflow uses [emojilib](https://www.npmjs.com/package/emojilib) and [emojiimages](https://www.npmjs.com/package/emojiimages).

For emojis in `emojilib`, if there is no image found in `emojiimages`, use a placeholder image instead.

Here is how the placeholder image generated:

```
convert -size 200x200 xc:#cacaca placeholder.png
```

To build the workflow:

- Clone the repo
- Install dependencies `yarn`
- Copy emoji images `cp -r node_modules/emojiimages/imgs workflow/`
- Generate emoji data file `go run scripts/generate_data.go`
- Build the workflow `make bundle`
