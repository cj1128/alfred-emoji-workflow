SHELL := /bin/bash

install:
	go install
.PHONY: install

bundle:
	go build -o workflow/emoji
	upx --brute workflow/emoji
	cd workflow && zip -r ../tmp/Emoji.alfredworkflow .
	rm -rf workflow/emoji
.PHONY: bundle
