# Makefile for AlphaArk Gin Project

.PHONY: help build run test clean deps lint

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  make build    - 构建项目"
	@echo "  make run      - 运行项目"
	@echo "  make test     - 运行测试"
	@echo "  make clean    - 清理构建文件"
	@echo "  make deps     - 安装依赖"
	@echo "  make lint     - 代码检查"

# 构建项目
build:
	@echo "构建项目..."
	go build -o bin/server cmd/server/main.go

# 运行项目
run:
	@echo "启动服务器..."
	go run cmd/server/main.go

# 运行测试
test:
	@echo "运行测试..."
	go test ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	go clean

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod tidy
	go mod download

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 生成文档
docs:
	@echo "生成文档..."
	swag init -g cmd/server/main.go

# 数据库迁移
migrate:
	@echo "数据库迁移..."
	# TODO: 添加数据库迁移命令

# Docker构建
docker-build:
	@echo "构建Docker镜像..."
	docker build -t alphaark-gin .

# Docker运行
docker-run:
	@echo "运行Docker容器..."
	docker run -p 8080:8080 alphaark-gin 