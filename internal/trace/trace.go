// Package trace 提供轻量埋点接口
package trace

import (
	"context"
)

// Span 表示一个追踪跨度
type Span interface {
	// AddEvent 添加事件
	AddEvent(name string, attributes map[string]interface{})

	// SetAttributes 设置属性
	SetAttributes(attributes map[string]interface{})

	// End 结束跨度
	End()
}

// Tracer 追踪器接口
type Tracer interface {
	// StartSpan 开始新的跨度
	StartSpan(ctx context.Context, name string) (context.Context, Span)
}

// noopTracer 空实现追踪器
type noopTracer struct{}

// noopSpan 空实现跨度
type noopSpan struct{}

// NewNoopTracer 创建空实现追踪器
func NewNoopTracer() Tracer {
	return &noopTracer{}
}

// StartSpan 实现Tracer接口
func (t *noopTracer) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	return ctx, &noopSpan{}
}

// AddEvent 实现Span接口
func (s *noopSpan) AddEvent(name string, attributes map[string]interface{}) {
	// 空实现
}

// SetAttributes 实现Span接口
func (s *noopSpan) SetAttributes(attributes map[string]interface{}) {
	// 空实现
}

// End 实现Span接口
func (s *noopSpan) End() {
	// 空实现
}

// 全局追踪器实例
var globalTracer Tracer = NewNoopTracer()

// SetGlobalTracer 设置全局追踪器
func SetGlobalTracer(tracer Tracer) {
	globalTracer = tracer
}

// StartSpan 使用全局追踪器开始跨度
func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	return globalTracer.StartSpan(ctx, name)
}
