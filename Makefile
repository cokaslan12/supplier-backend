build:
	@go build -o bin/api

run: build
	@./bin/api --listenAddr :2000

test:
	@go test -v ./...

