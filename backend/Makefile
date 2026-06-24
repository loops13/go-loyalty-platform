run:
	go run ./cmd/server

build:
	go build -o bin/goloyaltyplatform ./cmd/server

docker-build:
	docker build -t goloyaltyplatform:latest .

fmt:
	gofmt -s -w .

swagger:
	swag init -g ./cmd/server/main.go -o ./docs --parseInternal --parseDependency
