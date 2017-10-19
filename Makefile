SHELL := /bin/bash

install:
	go install
.PHONY: install

build:
	go build -o workflow/emoji
.PHONY: build

generate-data:
	go run scripts/generate_data.go
.PHONY: generate-data
