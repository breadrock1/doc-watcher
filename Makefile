BIN_LINTER := "${GOPATH}/bin/golangci-lint"
BIN_NATIVE_NOTIFIER := "./bin/doc-notifier"
BIN_MINIO_NOTIFIER := "./bin/doc-notifier-minio"
DOCKER_IMAGE := "doc-notifier:latest"

GIT_HASH := $(shell git log --format="%h" -n 1)

build:
	go build -v -o $(BIN_NATIVE_NOTIFIER) ./cmd/notifier
	go build -v -o $(BIN_MINIO_NOTIFIER) ./cmd/minio

run: build
	$(BIN_NOTIFIER) -c ./configs/config.toml

test:
	go test -race ./...

.PHONY: build run test