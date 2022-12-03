APP=go_rtsp_to_web
SERVER_FLAGS ?= -config config/config.json

P="\\033[34m[+]\\033[0m"

build:
	@echo "$(P) build"
	go build *.go

run:
	@echo "$(P) run"
	go run main.go

serve:
	@$(MAKE) server

server:
	@echo "$(P) server $(SERVER_FLAGS)"
	./${APP} $(SERVER_FLAGS)

test:
	@echo "$(P) test"
	bash test.curl
	bash test_multi.curl

lint:
	@echo "$(P) lint"
	go vet

.NOTPARALLEL:

.PHONY: build run server test lint