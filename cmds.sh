# Generate Swagger docs
swag init --parseDependency --parseInternal -g main.go --exclude infinirewards

swagger-markdown -i docs/swagger.json -o API.md

# Generate markdown docs
swag fmt

# Run tests
go test -p 1 -v ./api/tests/...

# Generate coverage report
go tool cover -html=coverage.out

# Run tests with coverage
go test -p 1 -coverprofile=coverage.out -v ./api/tests/...
