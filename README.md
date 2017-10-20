# Alfred Emoji Workflow

[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://mit-license.org/2016)

Alfred workflow to quick search emojis, after all, every one loves emoji 😎

## Generate emojis.go

1. install emojilib

```bash
git clone https://github.com/fate-lovely/emojilib.git tmp/
cp -r tmp/emojilib/imgs workflow
```

2. parse emojilib/emojis.json to get emojis.go

```bash
go run scripts/generate_data.go
```
