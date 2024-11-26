build:
	@go build -o bin/api

run: build
	@./bin/api --listenAddr :2000

seed:
	@go run scripts/seed.go

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 2000:2000 api	

test:
	@go test -v ./...