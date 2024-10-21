FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache gcc libc-dev make

WORKDIR /app

COPY . .

RUN go mod download && make build


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app .

RUN mkdir -p indexer && mkdir -p uploads

EXPOSE 2893

ENTRYPOINT [ "/app/bin/doc-notifier-minio", "-e" ]
