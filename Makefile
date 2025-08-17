# Makefile for BTC APIs project
# 提供常用的开发命令和任务

.PHONY: help build test clean lint format example

# 默认目标
.DEFAULT_GOAL := help

# 项目信息
PROJECT_NAME := btcapis
VERSION := 1.0.0
GO_VERSION := 1.21

# 颜色定义
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

help: ## 显示帮助信息
	@echo "$(GREEN)=== $(PROJECT_NAME) v$(VERSION) ===$(NC)"
	@echo "$(YELLOW)专注领域: 比特币区块链 API$(NC)"
	@echo "$(YELLOW)特色: 一次导入，所有功能立即可用$(NC)"
	@echo "$(YELLOW)可用的命令:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

build: ## 构建项目
	@echo "$(GREEN)构建BTC API项目...$(NC)"
	@go build -o bin/$(PROJECT_NAME) ./examples/usage_example.go

test: ## 运行所有测试
	@echo "$(GREEN)运行BTC API测试...$(NC)"
	@go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "$(GREEN)运行测试并生成覆盖率报告...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)覆盖率报告已生成: coverage.html$(NC)"

test-bench: ## 运行性能基准测试
	@echo "$(GREEN)运行性能基准测试...$(NC)"
	@go test -bench=. -benchmem ./...

lint: ## 运行代码检查
	@echo "$(GREEN)运行代码检查...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint 未安装，跳过代码检查$(NC)"; \
	fi

format: ## 格式化代码
	@echo "$(GREEN)格式化代码...$(NC)"
	@go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "$(YELLOW)goimports 未安装，跳过导入格式化$(NC)"; \
	fi

clean: ## 清理构建文件
	@echo "$(GREEN)清理构建文件...$(NC)"
	@rm -rf bin/
	@rm -rf coverage.out
	@rm -rf coverage.html
	@go clean -cache

deps: ## 下载和整理依赖
	@echo "$(GREEN)下载依赖...$(NC)"
	@go mod download
	@go mod tidy

deps-update: ## 更新依赖到最新版本
	@echo "$(GREEN)更新依赖...$(NC)"
	@go get -u ./...
	@go mod tidy

example: ## 运行使用示例
	@echo "$(GREEN)运行BTC API使用示例...$(NC)"
	@go run ./examples/usage_example.go

install-tools: ## 安装开发工具
	@echo "$(GREEN)安装开发工具...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)开发工具安装完成$(NC)"

check-go-version: ## 检查Go版本
	@echo "$(GREEN)检查Go版本...$(NC)"
	@go version
	@if [ "$$(go version | awk '{print $$3}' | sed 's/go//')" != "$(GO_VERSION)" ]; then \
		echo "$(YELLOW)警告: 推荐使用Go $(GO_VERSION) 或更高版本$(NC)"; \
	else \
		echo "$(GREEN)Go版本检查通过$(NC)"; \
	fi

dev-setup: ## 开发环境设置
	@echo "$(GREEN)设置BTC API开发环境...$(NC)"
	@make check-go-version
	@make deps
	@make install-tools
	@echo "$(GREEN)开发环境设置完成$(NC)"

all: clean deps format lint test build ## 执行完整的构建流程

# 显示项目信息
info: ## 显示项目信息
	@echo "$(GREEN)项目名称:$(NC) $(PROJECT_NAME)"
	@echo "$(GREEN)版本:$(NC) $(VERSION)"
	@echo "$(GREEN)Go版本要求:$(NC) $(GO_VERSION)+"
	@echo "$(GREEN)专注领域:$(NC) 比特币区块链"
	@echo "$(GREEN)导入方式:$(NC) import \"github.com/yourusername/btcapis\""
	@echo "$(GREEN)项目路径:$(NC) $(shell pwd)"
	@echo "$(GREEN)Go模块:$(NC) $(shell go list -m)"
