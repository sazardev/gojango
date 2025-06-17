# GoJango Makefile

.PHONY: help build run test clean format lint install-deps example-basic example-advanced

# Default target
help: ## Show this help
	@echo "GoJango - Django-inspired web framework for Go"
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build
build: ## Build the project
	@echo "🔨 Building GoJango..."
	go build -v ./...
	@echo "✅ Build complete!"

# Run tests
test: ## Run all tests
	@echo "🧪 Running tests..."
	go test -v ./...
	@echo "✅ Tests complete!"

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Benchmark tests
benchmark: ## Run benchmark tests
	@echo "⚡ Running benchmarks..."
	go test -bench=. ./...

# Format code
format: ## Format Go code
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted!"

# Lint code
lint: ## Lint Go code
	@echo "🔍 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi
	@echo "✅ Linting complete!"

# Install dependencies
install-deps: ## Install project dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	go clean
	rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"

# Run simple example
example-simple: ## Run the simple example
	@echo "🚀 Running simple example..."
	cd examples/simple && go run main.go

# Run basic example
example-basic: ## Run the basic example
	@echo "🚀 Running basic example..."
	cd examples/basic && go run main.go

# Run advanced example
example-advanced: ## Run the advanced example
	@echo "🚀 Running advanced example..."
	cd examples/advanced && go run main.go

# Development server with hot reload (requires air)
dev: ## Start development server with hot reload
	@echo "🔥 Starting development server with hot reload..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  'air' not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Initialize new GoJango project
init: ## Initialize a new GoJango project
	@echo "🎉 Initializing new GoJango project..."
	@read -p "Enter project name: " name; \
	mkdir -p $$name/{models,handlers,middleware,templates,static/{css,js}}; \
	cd $$name && go mod init $$name; \
	echo "module $$name\n\ngo 1.22\n\nrequire (\n\tgithub.com/tu-usuario/gojango v0.1.0\n)" > go.mod; \
	echo 'package main\n\nimport (\n\t"gojango"\n\t"gojango/models"\n)\n\ntype User struct {\n\tmodels.Model\n\tName  string `json:"name" db:"name,not_null"`\n\tEmail string `json:"email" db:"email,unique,not_null"`\n}\n\nfunc (u *User) TableName() string {\n\treturn "users"\n}\n\nfunc main() {\n\tapp := gojango.New()\n\tapp.AutoMigrate(&User{})\n\tapp.RegisterCRUD("/api/users", &User{})\n\tapp.GET("/", func(c *gojango.Context) error {\n\t\treturn c.JSON(map[string]string{"message": "¡Hola GoJango!"})\n\t})\n\tapp.Run(":8000")\n}' > main.go; \
	echo "✅ Project '$$name' created! Run 'cd $$name && make run' to start."

# Run the project
run: ## Run the main application
	@echo "🚀 Starting GoJango application..."
	go run main.go

# Create release
release: test ## Create a release build
	@echo "📦 Creating release..."
	@read -p "Enter version (e.g., v1.0.0): " version; \
	git tag $$version; \
	git push origin $$version; \
	echo "✅ Release $$version created!"

# Update dependencies
update: ## Update all dependencies
	@echo "🔄 Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies updated!"

# Generate documentation
docs: ## Generate documentation
	@echo "📚 Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "📖 Documentation server: http://localhost:6060/pkg/gojango/"; \
		godoc -http=:6060; \
	else \
		echo "⚠️  godoc not found. Installing..."; \
		go install golang.org/x/tools/cmd/godoc@latest; \
		echo "📖 Documentation server: http://localhost:6060/pkg/gojango/"; \
		godoc -http=:6060; \
	fi

# Docker build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t gojango-app .
	@echo "✅ Docker image built!"

# Docker run
docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8000:8000 gojango-app

# Show project status
status: ## Show project status
	@echo "📊 GoJango Project Status"
	@echo "========================"
	@echo "Go version: $(shell go version)"
	@echo "Module: $(shell go list -m)"
	@echo "Dependencies:"
	@go list -m all | grep -v "$(shell go list -m)" | head -10
	@echo ""
	@echo "Files:"
	@find . -name "*.go" -type f | wc -l | xargs echo "Go files:"
	@find . -name "*.go" -type f -exec wc -l {} + | tail -1 | awk '{print "Lines of code: " $$1}'
	@echo ""
	@echo "Tests:"
	@find . -name "*_test.go" -type f | wc -l | xargs echo "Test files:"

# Install dev tools
install-tools: ## Install development tools
	@echo "🛠️  Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/godoc@latest
	@echo "✅ Development tools installed!"

# Quick setup for new contributors
setup: install-deps install-tools ## Complete setup for new contributors
	@echo "🎉 Setup complete! You're ready to contribute to GoJango!"
	@echo ""
	@echo "Next steps:"
	@echo "  make test          # Run tests"
	@echo "  make example-basic # Try the basic example"
	@echo "  make dev          # Start development server"
