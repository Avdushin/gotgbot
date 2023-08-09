appname="tgbot"

Default:
	go build -o build/tgbot ./cmd/app/main.go
	./build/tgbot

build:
	go build -o build/tgbot ./cmd/app/main.go

run:
	go run ./cmd/app/main.go