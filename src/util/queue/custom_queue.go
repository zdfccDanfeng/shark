package queue

import (
	"github.com/eapache/queue"
	"runtime"
	"sync"
	"sync/atomic"
)

//MyQueue queue
type MyQueue struct {
	sync.Mutex
	popable *sync.Cond // Go的标准库中有一个类型叫条件变量：sync.Cond。这种类型与互斥锁和读写锁不同，它不是开箱即用的，它需要与互斥锁组合
	// ,类似于java里面的Condition
	buffer *queue.Queue
	closed bool
	count  int32
}

//New 创建
func New() *MyQueue {
	ch := &MyQueue{
		buffer: queue.New(),
	}
	ch.popable = sync.NewCond(&ch.Mutex)
	return ch
}

//Pop 取出队列,（阻塞模式）
func (q *MyQueue) Pop() (v interface{}) {
	c := q.popable // 获取condition
	buffer := q.buffer
	// 将竟态代码包裹在lock/unlock代码块里面，避免竞争访问
	q.Mutex.Lock()
	defer q.Mutex.Unlock() // 在程序返回之前完成锁释放

	for q.Len() == 0 && !q.closed {
		c.Wait() // 阻塞 释放锁许可
	}

	if q.closed { //已关闭
		return
	}

	if q.Len() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		atomic.AddInt32(&q.count, -1)
	}
	return
}

//试着取出队列（非阻塞模式）返回ok == false 表示空
func (q *MyQueue) TryPop() (v interface{}, ok bool) {
	buffer := q.buffer

	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if q.Len() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		atomic.AddInt32(&q.count, -1)
		ok = true
	} else if q.closed {
		ok = true
	}

	return
}

// 插入队列，非阻塞
func (q *MyQueue) Push(v interface{}) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.closed {
		q.buffer.Add(v)
		atomic.AddInt32(&q.count, 1)
		q.popable.Signal()
	}
}

// 获取队列长度
func (q *MyQueue) Len() int {
	return (int)(atomic.LoadInt32(&q.count))
}

// Close MyQueue
// After close, Pop will return nil without block, and TryPop will return v=nil, ok=True
func (q *MyQueue) Close() {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if !q.closed {
		q.closed = true
		atomic.StoreInt32(&q.count, 0)
		q.popable.Broadcast() //广播 notifyAll
	}
}

//Wait 等待队列消费完成
func (q *MyQueue) Wait() {
	for {
		if q.closed || q.Len() == 0 {
			break
		}

		runtime.Gosched() //出让时间片
	}
}
