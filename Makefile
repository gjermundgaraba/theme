.PHONY: build link serve test

build:
	go run ./cli build

link:
	go run ./cli link

serve:
	go run ./cli serve

test:
	go test ./...
