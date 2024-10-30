build:
	@CGO_ENABLED=0 GOOS=linux go build -mod=mod -o queuer -ldflags "-s -w" .
docker:
	@docker compose up --build -d
lint:
	@go mod tidy
	@gofumpt -d -w .
	@golangci-lint run

.PHONY: build lint docker
