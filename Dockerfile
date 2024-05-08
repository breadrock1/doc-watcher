FROM alpine:3.19.1

WORKDIR /app
COPY . .

RUN mkdir -p ./indexer/uploads ./indexer/watcher ./indexer/unrecognized

RUN apk add --no-cache curl

EXPOSE 2893

CMD ["/app/doc-notifier", "-ej"]
