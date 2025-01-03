FROM golang:1.22-alpine AS builder

RUN apk update && apk add --no-cache gcc libc-dev make

WORKDIR /app

COPY . .

RUN go mod download && make build


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app .

ENTRYPOINT [ "/app/bin/doc-watcher", "-e" ]

EXPOSE 2893
