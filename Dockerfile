FROM golang:1.21

WORKDIR /app

COPY . .

RUN rm -rf .env && go mod download
RUN go build -o fs-notifier .

EXPOSE 2550

CMD ["./fs-notifier"]
