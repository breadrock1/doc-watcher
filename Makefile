BIN_LINTER := "${GOPATH}/bin/golangci-lint"
BIN_NATIVE_NOTIFIER := "./bin/doc-notifier"
BIN_MINIO_NOTIFIER := "./bin/doc-notifier-minio"
BIN_DIR_EXPORTER := "./bin/doc-exporter"
DOCKER_IMAGE := "doc-notifier:latest"

GIT_HASH := $(shell git log --format="%h" -n 1)

build:
	go build -v -o $(BIN_NATIVE_NOTIFIER) ./cmd/notifier
	go build -v -o $(BIN_MINIO_NOTIFIER) ./cmd/minio
	go build -v -o $(BIN_DIR_EXPORTER) ./cmd/exporter

run-native: build
	$(BIN_NATIVE_NOTIFIER) -c ./configs/config.toml

run-minio: build
	$(BIN_MINIO_NOTIFIER) -c ./configs/config.toml

run-exporter: build
	$(BIN_DIR_EXPORTER) -e

test:
	go test -race ./...

.PHONY: build run test