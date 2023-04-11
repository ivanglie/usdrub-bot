FROM golang:1.19-alpine AS builder
WORKDIR /usr/src/usdrub-bot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -v -o usdrub-bot ./cmd/app

FROM --platform=$BUILDPLATFORM alpine:3.17.0
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
COPY --from=builder /usr/src/usdrub-bot/usdrub-bot /usr/local/bin/usdrub-bot
EXPOSE 18001
CMD ["usdrub-bot"]