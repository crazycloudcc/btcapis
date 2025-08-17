//go:build tools
// +build tools

// 此文件用于管理开发依赖，确保版本一致性
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "golang.org/x/tools/cmd/goimports"
)
