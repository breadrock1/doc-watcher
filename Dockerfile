FROM golang:1.21-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download && go build -o ./doc-notifier ./cmd/notifier

FROM golang:1.21-alpine

WORKDIR /app

COPY --from=builder /app/doc-notifier .
RUN mkdir indexer && mkdir upload

CMD ["/app/doc-notifier", "-e"]

EXPOSE 2893