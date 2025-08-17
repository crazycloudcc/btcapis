// Package chain 错误定义
package chain

import (
	"errors"
	"fmt"
)

// 预定义错误类型
var (
	// ErrNoBackends 表示没有可用的后端
	ErrNoBackends = errors.New("no backends available")

	// ErrBackendUnavailable 表示后端服务不可用
	ErrBackendUnavailable = errors.New("backend unavailable")

	// ErrBackendTimeout 表示后端响应超时
	ErrBackendTimeout = errors.New("backend timeout")

	// ErrBackendUnhealthy 表示后端不健康
	ErrBackendUnhealthy = errors.New("backend unhealthy")

	// ErrUnsupportedOperation 表示不支持的操作
	ErrUnsupportedOperation = errors.New("unsupported operation")
)

// RoutingError 表示路由错误
type RoutingError struct {
	Operation string
	Err       error
}

func (e *RoutingError) Error() string {
	return fmt.Sprintf("routing error for operation %s: %v", e.Operation, e.Err)
}

func (e *RoutingError) Unwrap() error {
	return e.Err
}

// NewRoutingError 创建新的路由错误
func NewRoutingError(operation string, err error) *RoutingError {
	return &RoutingError{
		Operation: operation,
		Err:       err,
	}
}
