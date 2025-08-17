// Package chain 多后端聚合、故障转移、负载策略
package chain

import (
	"context"
	"sync"
	"time"

	"github.com/crazycloudcc/btcapis/types"
)

// Router 定义后端路由策略
type Router struct {
	backends []Backend
	mu       sync.RWMutex

	// 路由策略
	strategy RoutingStrategy

	// 健康检查配置
	healthCheckInterval time.Duration
	healthCheckTimeout  time.Duration
}

// RoutingStrategy 定义路由策略接口
type RoutingStrategy interface {
	// SelectBackend 选择后端
	SelectBackend(ctx context.Context, backends []Backend, operation string) (Backend, error)

	// Name 获取策略名称
	Name() string
}

// NewRouter 创建新的路由器
func NewRouter(strategy RoutingStrategy) *Router {
	return &Router{
		strategy:            strategy,
		healthCheckInterval: 30 * time.Second,
		healthCheckTimeout:  5 * time.Second,
	}
}

// AddBackend 添加后端
func (r *Router) AddBackend(backend Backend) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.backends = append(r.backends, backend)
}

// RemoveBackend 移除后端
func (r *Router) RemoveBackend(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, backend := range r.backends {
		if backend.Name() == name {
			r.backends = append(r.backends[:i], r.backends[i+1:]...)
			break
		}
	}
}

// GetBackends 获取所有后端
func (r *Router) GetBackends() []Backend {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Backend, len(r.backends))
	copy(result, r.backends)
	return result
}

// SelectBackend 选择后端
func (r *Router) SelectBackend(ctx context.Context, operation string) (Backend, error) {
	r.mu.RLock()
	backends := r.GetBackends()
	r.mu.RUnlock()

	if len(backends) == 0 {
		return nil, types.ErrBackendUnavailable
	}

	return r.strategy.SelectBackend(ctx, backends, operation)
}

// StartHealthCheck 启动健康检查
func (r *Router) StartHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(r.healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.performHealthCheck(ctx)
		}
	}
}

// performHealthCheck 执行健康检查
func (r *Router) performHealthCheck(ctx context.Context) {
	backends := r.GetBackends()

	for _, backend := range backends {
		go func(b Backend) {
			ctx, cancel := context.WithTimeout(ctx, r.healthCheckTimeout)
			defer cancel()

			// 异步健康检查
			_ = b.IsHealthy(ctx)
		}(backend)
	}
}

// 预定义路由策略
type (
	// PriorityRouting 优先级路由策略
	PriorityRouting struct {
		priorities map[string]int
	}

	// RoundRobinRouting 轮询路由策略
	RoundRobinRouting struct {
		current int
		mu      sync.Mutex
	}

	// LoadBalancedRouting 负载均衡路由策略
	LoadBalancedRouting struct {
		mu sync.RWMutex
	}
)
