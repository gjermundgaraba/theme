.PHONY: link serve test

link:
	go run ./cli link

serve:
	go run ./cli serve

test:
	go test ./...
