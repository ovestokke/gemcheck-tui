.PHONY: build run test clean

build:
	go build -o gemcheck ./cmd/gemcheck

run: build
	./gemcheck

test:
	go test ./...

clean:
	rm -f gemcheck
