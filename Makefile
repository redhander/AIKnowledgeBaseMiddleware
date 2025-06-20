# AI Knowledge Base Middleware Makefile

.PHONY: help build run dev clean test install-air

# 默认目标
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the server"
	@echo "  dev        - Run with hot reload (requires air)"
	@echo "  install-air - Install air for hot reload"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"

# 构建应用
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go
	go build -o bin/cli cmd/cli/main.go
	@echo "Build completed!"

# 运行服务器
run:
	@echo "Starting server..."
	go run cmd/server/main.go

# 热加载开发模式
dev:
	@echo "Starting development server with hot reload..."
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "Air not found. Installing air..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

# 安装air工具
install-air:
	@echo "Installing air for hot reload..."
	go install github.com/air-verse/air@latest
	@echo "Air installed successfully!"

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean
	@echo "Clean completed!"

# 运行测试
test:
	@echo "Running tests..."
	go test ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# 构建Docker镜像
docker-build:
	@echo "Building Docker image..."
	docker build -t ai-knowledge-base-middleware .

# 运行Docker容器
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 ai-knowledge-base-middleware
