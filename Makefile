BOT_TOKEN=$(shell cat bottoken.secret)

tests:
	go test -v -cover -race ./...

run:
	go run -race ./cmd/app/ --bottoken=$(BOT_TOKEN)

dc:
	BOT_TOKEN=$(BOT_TOKEN) docker compose up -d