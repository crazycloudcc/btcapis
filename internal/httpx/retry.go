// Package httpx 重试机制
package httpx

import (
	"net/http"
	"time"
)

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxRetries int
	Delay      time.Duration
	Backoff    BackoffStrategy
	Retryable  RetryableCondition
}

// BackoffStrategy 退避策略
type BackoffStrategy interface {
	NextDelay(attempt int, delay time.Duration) time.Duration
}

// RetryableCondition 可重试条件
type RetryableCondition interface {
	ShouldRetry(resp *http.Response, err error) bool
}

// ExponentialBackoff 指数退避策略
type ExponentialBackoff struct {
	Multiplier float64
	MaxDelay   time.Duration
}

// NextDelay 计算下一次延迟
func (e *ExponentialBackoff) NextDelay(attempt int, delay time.Duration) time.Duration {
	if attempt == 0 {
		return delay
	}

	newDelay := time.Duration(float64(delay) * e.Multiplier)
	if newDelay > e.MaxDelay {
		newDelay = e.MaxDelay
	}

	return newDelay
}

// DefaultRetryableCondition 默认可重试条件
type DefaultRetryableCondition struct{}

// ShouldRetry 判断是否应该重试
func (d *DefaultRetryableCondition) ShouldRetry(resp *http.Response, err error) bool {
	// 网络错误可以重试
	if err != nil {
		return true
	}

	// 5xx错误可以重试
	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return true
	}

	// 429 Too Many Requests可以重试
	if resp.StatusCode == 429 {
		return true
	}

	return false
}

// NewRetryPolicy 创建重试策略
func NewRetryPolicy(maxRetries int, delay time.Duration) *RetryPolicy {
	return &RetryPolicy{
		MaxRetries: maxRetries,
		Delay:      delay,
		Backoff:    &ExponentialBackoff{Multiplier: 2.0, MaxDelay: 30 * time.Second},
		Retryable:  &DefaultRetryableCondition{},
	}
}
