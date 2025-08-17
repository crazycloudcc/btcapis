// Package httpx 速率限制器
package httpx

import (
	"context"
	"sync"
	"time"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	rate       int       // 每秒请求数
	burst      int       // 突发请求数
	tokens     int       // 当前令牌数
	lastRefill time.Time // 上次补充时间
	mu         sync.Mutex
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(rate, burst int) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		burst:      burst,
		tokens:     burst,
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许请求
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill)

	// 计算需要补充的令牌数
	tokensToAdd := int(elapsed.Seconds() * float64(r.rate))
	if tokensToAdd > 0 {
		r.tokens = min(r.burst, r.tokens+tokensToAdd)
		r.lastRefill = now
	}

	if r.tokens > 0 {
		r.tokens--
		return true
	}

	return false
}

// Wait 等待直到允许请求
func (r *RateLimiter) Wait(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if r.Allow() {
				return nil
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
