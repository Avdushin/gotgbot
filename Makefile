appname="tgbot"

Default:
	go run ./cmd/app/main.go

Build:
	go build -o tgbot ./cmd/app/main.go