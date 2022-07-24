FROM golang:1.17-alpine

WORKDIR /usr/src/usdrub-bot/

COPY . .

RUN go build -v -o /usr/local/bin/usdrub-bot ./cmd/app/

CMD ["usdrub-bot"]