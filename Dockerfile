FROM golang:1.21-alpine as builder

RUN apk update && apk add --no-cache gcc libc-dev make

WORKDIR /app

COPY . .

RUN go mod download && make build

FROM golang:1.21-alpine

WORKDIR /app

COPY --from=builder /app .
RUN mkdir -p indexer && mkdir -p uploads

CMD ["/app/bin/doc-notifier-minio", "-e"]

EXPOSE 2893