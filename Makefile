.PHONY: all build run tui test clean docker-build docker-run

BINARY_NAME=banana-ssh
PORT?=2222
HOST?=0.0.0.0

all: build test

build:
	go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/banana-ssh

run: build
	./$(BINARY_NAME) -port $(PORT) -host $(HOST)

tui: build
	./$(BINARY_NAME) -tui

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)
	rm -rf .ssh

docker-build:
	docker build -t banana-ssh:latest .

docker-run:
	docker run -it --rm -p $(PORT):2222 -v banana_ssh_data:/data banana-ssh:latest
