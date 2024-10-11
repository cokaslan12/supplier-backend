build:
	@go build -o bin/api

run: build
	@./bin/api --listenAddr :2000

