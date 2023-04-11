.SILENT:

default: build

build:
	go build -o .bin/crypto-currency cmd/main.go

test: build
	./.bin/crypto-currency | jq .

prod: build
	\cp -f .bin/crypto-currency ~/.config/waybar/scripts/
