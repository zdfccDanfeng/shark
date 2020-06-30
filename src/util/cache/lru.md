![lru](https://geektutu.com/post/geecache-day1/lru.jpg)

- 绿色的是字典(map)，存储键和值的映射关系。这样根据某个键(key)查找对应的值(value)的复杂是O(1)，在字典中插入一条记录的复杂度也是O(1)。
- 红色的是双向链表(double linked list)实现的队列。将所有的值放到双向链表中，这样，当访问到某个值时，将其移动到队尾的复杂度是O(1)，
在队尾新增一条记录以及删除一条记录的复杂度均为O(1)。

互斥锁<br>
  &nbsp;&nbsp;&nbsp;&nbsp;多个协程(goroutine)同时读写同一个变量，在并发度较高的情况下，会发生冲突。确保一次只有一个协程(goroutine)可以访问该变量以避免冲突，这称之为互斥，互斥锁可以解决这个问题。
  // sync.Mutex 是一个互斥锁，可以由不同的协程加锁和解锁
  
  
  
- Group结构
```cfml
                是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶

// cache 雏形结构

testLru/
    |--lru/
        |--lru.go  // lru 缓存淘汰策略
    |--byteview.go // 缓存值的抽象与封装
    |--cache.go    // 并发控制
    |--geecache.go // 负责与外部交互，控制缓存存储和获取的主流程
```

https://github.com/geektutu/7days-golang/blob/master/gee-cache/day1-lru/geecache/lru/lru.go

https://go.wuhaolin.cn/gopl/ch4/ch4-02.html // go程序设计

