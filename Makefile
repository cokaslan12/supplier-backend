build:
	@go build -o bin/api

run: build
	@./bin/api --listenAddr :2000

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...