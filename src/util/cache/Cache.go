package cache

import (
	"github.com/shark/src/util/cache/lru"
	"sync"
)

// 支持并发访问读写的缓存结构 --- 在原来的lru的基础上增加了锁结构来实现并发读写控制
// 实例化 lru，封装 get 和 add 方法，并添加互斥锁 mu
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 惰性操作
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
