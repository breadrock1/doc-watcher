FROM golang:1.21-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download && make build

FROM golang:1.21-alpine

WORKDIR /app

COPY --from=builder /app/doc-notifier .
RUN mkdir indexer && mkdir uploads

CMD ["/app/bin/doc-notifier-minio", "-e"]

EXPOSE 2893