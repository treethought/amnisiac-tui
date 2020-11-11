help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

SHELL := /bin/bash

mpv:

	mpv --idle=once --no-terminal --input-ipc-server=/tmp/mpvsocket --no-video &
	# mpv --idle --input-ipc-server=/tmp/mpvsocket --no-terminal &

run: mpv ## run locally outside of docker
	source .env && go run main.go
	killall mpv


compile: ## build go binary
	go build -v -o amnisiac .

build: ## build image locally
	docker build -t didi-bots .


test: emulator ## run tests
	go test -i
	source .env && gotestsum --format testname -- -v -failfast -coverprofile coverage.out ./...
	docker kill ds-emulator

cover: ## view test coverage
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

