docker:
	@docker compose up --build -d
lint:
	@go mod tidy
	@gofumpt -d -w .
	@golangci-lint run

.PHONY: lint docker
