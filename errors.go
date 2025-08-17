// Package btcapis 错误定义
package btcapis

import (
	"errors"
	"fmt"
)

// 预定义错误类型
var (
	// ErrNotFound 表示请求的资源未找到
	ErrNotFound = errors.New("resource not found")

	// ErrTimeout 表示操作超时
	ErrTimeout = errors.New("operation timeout")

	// ErrBackendUnavailable 表示后端服务不可用
	ErrBackendUnavailable = errors.New("backend unavailable")

	// ErrInvalidInput 表示输入参数无效
	ErrInvalidInput = errors.New("invalid input")

	// ErrNetworkMismatch 表示网络类型不匹配
	ErrNetworkMismatch = errors.New("network mismatch")

	// ErrUnsupportedOperation 表示不支持的操作
	ErrUnsupportedOperation = errors.New("unsupported operation")
)

// BackendError 表示后端服务错误
type BackendError struct {
	Backend string
	Err     error
}

func (e *BackendError) Error() string {
	return fmt.Sprintf("backend %s error: %v", e.Backend, e.Err)
}

func (e *BackendError) Unwrap() error {
	return e.Err
}

// NewBackendError 创建新的后端错误
func NewBackendError(backend string, err error) *BackendError {
	return &BackendError{
		Backend: backend,
		Err:     err,
	}
}

// ValidationError 表示验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field %s: %s", e.Field, e.Message)
}

// NewValidationError 创建新的验证错误
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
