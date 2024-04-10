lint:
	golangci-lint run --timeout=10m -v ./...


swag-api:
	swag init --dir ./api --output ./api/docs --parseDependency --parseInternal
