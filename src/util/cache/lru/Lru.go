package lru

import "container/list"

// see https://geektutu.com/post/geecache-day1.html#1-FIFO-LFU-LRU-%E7%AE%97%E6%B3%95%E7%AE%80%E4%BB%8B

// LRU 算法最核心的 2 个数据结构  Good Design !!!
//
//绿色的是字典(map)，存储键和值的映射关系。这样根据某个键(key)查找对应的值(value)的复杂是O(1)，在字典中插入一条记录的复杂度也是O(1)。
//红色的是双向链表(double linked list)实现的队列。将所有的值放到双向链表中，这样，当访问到某个值时，将其移动到队尾的复杂度是O(1)，
// 在队尾新增一条记录以及删除一条记录的复杂度均为O(1)。

// Cache is a LRU cache. It is not safe for concurrent access.
// 保证队列首部是最近最少使用的元素
type Cache struct {
	maxBytes int64 // 允许使用的最大内存
	nbytes   int64 // 当前已使用的内存
	ll       *list.List
	// 这里直接使用go语言标准库实现的双向List
	cache map[string]*list.Element // 键是字符串，值是双向链表中对应节点的指针
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可以为 nil
}

// 键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，
// 需要用 key 从字典中删除对应的映射
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
// 为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// New is the Constructor of Cache
// 方便实例化 Cache，实现 New() 函数：
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
// 查找主要有 2 个步骤，第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 将元素移动到队列尾部
		// c.ll.MoveToFront(ele)，即将链表中的节点 ele 移动到队尾（双向链表作为队列，队首队尾是相对的
		c.ll.MoveToFront(ele)    //头插入法，插入到链表头部位置
		kv := ele.Value.(*entry) // go泛型
		return kv.value, true
	}
	return
}

// 执行删除操作  淘汰在队列里最近最少访问的元素 ---》 也就是队列尾部的元素
// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	// 定位到队列的首部
	ele := c.ll.Back() // 取到队尾节点，从链表中删除
	if ele != nil {
		c.ll.Remove(ele)
		// 类型强转
		kv := ele.Value.(*entry)
		// 从字典里面删除掉这一个映射关系
		delete(c.cache, kv.key)
		// 更新当前使用的内存数量
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) // 如果回调函数 OnEvicted 不为 nil，则调用回调函数
		}
	}
}

// 增加元素移动到队列队头尾部
// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	// 先检测要添加的元素是否在缓存中
	if ele, ok := c.cache[key]; ok {
		// 如果在缓存中，则将其移动到队列头部，并进行缓存更新操作 【【【值覆盖操作】】】
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
		return
	}
	// 如果不在缓存中，则进行添加操作，并将元素移动到队列头部
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	// 更新使用的内存情况
	c.nbytes += int64(len(key)) + int64(value.Len())
	// 检测内存容量。如果容量超过了最大的内存限制，则移除最近最久没有访问的元素
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
