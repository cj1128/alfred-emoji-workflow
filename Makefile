SHELL := /bin/bash

install:
	go install
.PHONY: install

build:
	go build -o workflow/emoji

bundle: build
	upx --brute workflow/emoji
	cd workflow && zip -r ../tmp/Emoji.alfredworkflow .
.PHONY: bundle
