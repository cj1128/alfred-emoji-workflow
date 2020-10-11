SHELL := /bin/bash

build:
	go build -o workflow/emoji
.PHONY: build

bundle: build
	upx --brute workflow/emoji
	cd workflow && zip -r ../Emoji.alfredworkflow .
.PHONY: bundle
