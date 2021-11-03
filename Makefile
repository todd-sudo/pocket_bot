.PHONY:

build:
	go build -o .bin/bot cmd/bot/main.go

run: build
	go run cmd/bot/main.go