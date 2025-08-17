// Package httpx HTTP中间件
package httpx

import (
	"net/http"
	"time"
)

// Middleware HTTP中间件接口
type Middleware interface {
	Wrap(next http.RoundTripper) http.RoundTripper
}

// MiddlewareFunc 中间件函数类型
type MiddlewareFunc func(next http.RoundTripper) http.RoundTripper

// Wrap 实现Middleware接口
func (f MiddlewareFunc) Wrap(next http.RoundTripper) http.RoundTripper {
	return f(next)
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	Logger Logger
}

// Logger 日志接口
type Logger interface {
	Log(level string, msg string, fields map[string]interface{})
}

// Wrap 包装HTTP传输器
func (m *LoggingMiddleware) Wrap(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		start := time.Now()

		// 记录请求
		m.Logger.Log("info", "HTTP request started", map[string]interface{}{
			"method": req.Method,
			"url":    req.URL.String(),
		})

		// 执行请求
		resp, err := next.RoundTrip(req)

		// 记录响应
		duration := time.Since(start)
		fields := map[string]interface{}{
			"method":   req.Method,
			"url":      req.URL.String(),
			"duration": duration.String(),
		}

		if err != nil {
			m.Logger.Log("error", "HTTP request failed", fields)
		} else {
			fields["status"] = resp.StatusCode
			m.Logger.Log("info", "HTTP request completed", fields)
		}

		return resp, err
	})
}

// RoundTripperFunc 函数类型HTTP传输器
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip 实现http.RoundTripper接口
func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Chain 中间件链
type Chain struct {
	middlewares []Middleware
}

// NewChain 创建新的中间件链
func NewChain(middlewares ...Middleware) *Chain {
	return &Chain{
		middlewares: middlewares,
	}
}

// Wrap 包装HTTP传输器
func (c *Chain) Wrap(next http.RoundTripper) http.RoundTripper {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		next = c.middlewares[i].Wrap(next)
	}
	return next
}
