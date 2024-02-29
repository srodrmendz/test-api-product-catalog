test:
	go test -cover ./...

lint:
	golangci-lint run --fast ./...

deps:
	go mod tidy

docs:
	@swag init --parseDependency -g app.go 