FROM golang:1.19-alpine

RUN apk add --no-cache tzdata

ENV TZ=Europe/Moscow

WORKDIR /usr/src/usdrub-bot/

COPY . .

RUN go build -v -o /usr/local/bin/usdrub-bot ./cmd/app/

CMD ["usdrub-bot"]