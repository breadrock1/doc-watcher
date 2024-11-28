BIN_LINTER := "${GOPATH}/bin/golangci-lint"
BIN_MINIO_NOTIFIER := "./bin/doc-watcher"

GIT_HASH := $(shell git log --format="%h" -n 1)

build:
	go build -v -o $(BIN_MINIO_NOTIFIER) ./cmd/minio

run: build
	$(BIN_MINIO_NOTIFIER) -c ./configs/production.toml

test:
	go test -race ./tests/...

.PHONY: build run test