.PHONY: help build test lint clean deps fmt vet coverage integration

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  build        - 构建项目"
	@echo "  test         - 运行单元测试"
	@echo "  lint         - 运行代码质量检查"
	@echo "  clean        - 清理构建产物"
	@echo "  deps         - 下载依赖"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 运行go vet"
	@echo "  coverage     - 生成测试覆盖率报告"
	@echo "  integration  - 运行集成测试"

# 构建项目
build:
	go build -o bin/btcapis ./examples/basic

# 运行单元测试
test:
	go test -v ./...

# 运行代码质量检查
lint:
	golangci-lint run

# 清理构建产物
clean:
	rm -rf bin/
	go clean -cache -testcache

# 下载依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...
	goimports -w .

# 运行go vet
vet:
	go vet ./...

# 生成测试覆盖率报告
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 运行集成测试
integration:
	go test -tags=integration -v ./test/...

# 安装开发工具
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# 预提交检查
pre-commit: fmt lint test
