FROM golang:1.19-alpine AS builder
WORKDIR /usr/src/usdrub-bot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o usdrub-bot ./cmd/app

FROM scratch
COPY --from=builder /usr/src/usdrub-bot/usdrub-bot /usr/local/bin/usdrub-bot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Europe/Moscow
CMD ["usdrub-bot"]