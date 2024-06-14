BIN_NOTIFIER := "./bin/doc-notifier"
DOCKER_IMAGE="doc-notifier:latest"

GIT_HASH := $(shell git log --format="%h" -n 1)

build:
	go build -v -o $(BIN_NOTIFIER) ./cmd/notifier

run: build
	$(BIN_NOTIFIER) -c ./configs/config.toml

test:
	go test -race ./...

.PHONY: build run test