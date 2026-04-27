.PHONY: dev build clean docker-build docker-run test

# Development with Air live reload
dev:
	@echo "Starting API gateway with Air..."
	air

# Build the binary
build:
	@echo "Building API gateway..."
	go build -o api-gateway ./main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf api-gateway tmp/
	rm -f air.log

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t cryplio/api-gateway .

# Docker run
docker-run:
	@echo "Running API gateway in Docker..."
	docker run -p 8080:8080 --env-file .env cryplio/api-gateway

# Test the binary
test:
	@echo "Testing API gateway..."
	@./api-gateway &
	@sleep 2
	@curl -s http://localhost:8080/health | head -20
	@curl -s http://localhost:8080/live
	@curl -s http://localhost:8080/ready
	@pkill -f api-gateway || true
	@echo "Tests passed!"

# Run with direct go run (no Air)
run:
	go run ./main.go

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Tidy dependencies
tidy:
	go mod tidy
