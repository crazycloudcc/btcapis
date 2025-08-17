// Package bytespool 提供字节切片对象池
package bytespool

import (
	"sync"
)

// Pool 字节切片对象池
type Pool struct {
	pool sync.Pool
}

// NewPool 创建新的字节池
func NewPool() *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024) // 默认容量1024字节
			},
		},
	}
}

// Get 从池中获取字节切片
func (p *Pool) Get() []byte {
	return p.pool.Get().([]byte)
}

// Put 将字节切片放回池中
func (p *Pool) Put(b []byte) {
	// 重置切片长度，保留容量
	b = b[:0]
	p.pool.Put(b)
}

// GetWithCapacity 获取指定容量的字节切片
func (p *Pool) GetWithCapacity(capacity int) []byte {
	b := p.Get()
	if cap(b) < capacity {
		// 如果池中的切片容量不够，创建新的
		return make([]byte, 0, capacity)
	}
	return b
}

// 全局字节池实例
var globalPool = NewPool()

// Get 使用全局池获取字节切片
func Get() []byte {
	return globalPool.Get()
}

// Put 使用全局池放回字节切片
func Put(b []byte) {
	globalPool.Put(b)
}

// GetWithCapacity 使用全局池获取指定容量的字节切片
func GetWithCapacity(capacity int) []byte {
	return globalPool.GetWithCapacity(capacity)
}
