.PHONY: all build run tui test clean docker-build docker-run release-build

BINARY_NAME=banana
PORT?=2222
HOST?=0.0.0.0

VERSION?=1.0.0
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "head")
DATE?=$(shell date -u +%Y-%m-%d)
LDFLAGS=-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)

all: build test

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/banana-ssh

run: build
	./$(BINARY_NAME) serve -port $(PORT) -host $(HOST)

tui: build
	./$(BINARY_NAME)

test:
	go test -v ./...

release-build:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/banana-linux-amd64 ./cmd/banana-ssh
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/banana-linux-arm64 ./cmd/banana-ssh
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/banana-darwin-amd64 ./cmd/banana-ssh
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/banana-darwin-arm64 ./cmd/banana-ssh
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/banana-windows-amd64.exe ./cmd/banana-ssh
	cd dist && sha256sum banana-* > SHA256SUMS.txt

clean:
	rm -f $(BINARY_NAME) banana-ssh
	rm -rf dist .ssh

docker-build:
	docker build -t banana-ssh:latest .

docker-run:
	docker run -it --rm -p $(PORT):2222 -v banana_ssh_data:/data banana-ssh:latest
