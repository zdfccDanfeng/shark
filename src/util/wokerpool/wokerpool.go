package wokerpool

import (
	"context"
	"fmt"
	"runtime"
	"scaffold/src/util/queue"
	"sync/atomic"
	"testing"
	"time"
)

// New new workpool and set the max number of concurrencies
func New(max int) *WorkPool { // 注册工作池，并设置最大并发数
	if max < 1 {
		max = 1
	}
	// 创建Worker_pool
	p := &WorkPool{
		task:         make(chan TaskHandler, 2*max),
		errChan:      make(chan error, 1),
		waitingQueue: queue.New(), // 线程安全的双buffer阻塞队列
	}
	// 开启任务循环逻辑♻️
	go p.loop(max)
	return p
}

// SetTimeout Setting timeout time
func (p *WorkPool) SetTimeout(timeout time.Duration) { // 设置超时时间
	p.timeout = timeout
}

// Do Add to the workpool and return immediately
func (p *WorkPool) Do(fn TaskHandler) { // 添加到工作池，并立即返回
	if p.IsClosed() { // 已关闭
		return
	}
	p.waitingQueue.Push(fn)
	// p.task <- fn
}

// DoWait Add to the workpool and wait for execution to complete before returning
func (p *WorkPool) DoWait(task TaskHandler) { // 添加到工作池，并等待执行完成之后再返回
	if p.IsClosed() { // closed
		return
	}

	doneChan := make(chan struct{})
	p.waitingQueue.Push(TaskHandler(func() error {
		defer close(doneChan)
		return task()
	}))
	<-doneChan
}

// Wait Waiting for the worker thread to finish executing
func (p *WorkPool) Wait() error { // 等待工作线程执行结束
	p.waitingQueue.Wait()  // 等待队列结束
	p.waitingQueue.Close() //
	p.waitTask()           // wait que down
	close(p.task)          // 关闭任务处理Channel
	p.wg.Wait()            // 等待结束
	select {
	case err := <-p.errChan:
		return err
	default:
		return nil
	}
}

// IsDone Determine whether it is complete (non-blocking)
func (p *WorkPool) IsDone() bool { // 判断是否完成 (非阻塞)
	if p == nil || p.task == nil {
		return true
	}

	return p.waitingQueue.Len() == 0 && len(p.task) == 0
}

// IsClosed Has it been closed?
func (p *WorkPool) IsClosed() bool { // 是否已经关闭
	if atomic.LoadInt32(&p.closed) == 1 { // closed
		return true
	}
	return false
}

func (p *WorkPool) startQueue() {
	p.isQueTask = 1
	// 死循环 从任务队列里面拉取任务放到worker_pool的 goroutine 里面进行执行
	for {
		tmp := p.waitingQueue.Pop()
		if p.IsClosed() { // closed
			p.waitingQueue.Close()
			break
		}
		if tmp != nil {
			// 如果从Queue里面取出的任务不为空，则执行该任务
			// 参看：in.(type) ，将interface类型转换成具体的type类型
			fn := tmp.(TaskHandler)
			if fn != nil {
				p.task <- fn // 放入任务到worker_pool 的task_handler chan
			}
		} else {
			break
		}
	}
	atomic.StoreInt32(&p.isQueTask, 0)
}

// 阻塞直到任务执行完成
func (p *WorkPool) waitTask() {
	for {
		runtime.Gosched() // 出让时间片
		if p.IsDone() {
			if atomic.LoadInt32(&p.isQueTask) == 0 {
				break
			}
		}
	}
}

// 限制goroutine 的数量，防止goroutine 爆炸
func (p *WorkPool) loop(maxWorkersCount int) {
	go p.startQueue() // Startup queue , 启动队列

	p.wg.Add(maxWorkersCount) // Maximum number of work cycles,最大的工作协程数 ，类似于java的count_down_latch,开启多个任务
	// Start Max workers, 启动max个worker
	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			defer p.wg.Done() // countDownLatch的#countDown方法
			// worker 开始干活
			for wt := range p.task {
				if wt == nil || atomic.LoadInt32(&p.closed) == 1 { // returns immediately,有err 立即返回
					continue // It needs to be consumed before returning.需要先消费完了之后再返回，
				}
				// 关闭信号
				closed := make(chan struct{}, 1)
				// Set timeout, priority task timeout.有设置超时,优先task 的超时
				if p.timeout > 0 {
					// 每一个 context.Context 都会从最顶层的 Goroutine 一层一层传递到最下层。
					// context.Context 可以在上层 Goroutine 执行出现错误时，将信号及时同步给下层。
					// @see https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/
					ct, cancel := context.WithTimeout(context.Background(), p.timeout)
					// 后台监听任务完成状态。。。。
					go func() {
						select {
						case <-ct.Done():
							p.errChan <- ct.Err()
							// if atomic.LoadInt32(&p.closed) != 1 {
							// mylog.Error(ct.Err())
							atomic.StoreInt32(&p.closed, 1)
							cancel()
						case <-closed:
						}
					}()
				}
				// 执行任务逻辑
				err := wt() // Points of Execution.真正执行的点
				close(closed)
				if err != nil {
					select {
					// 如果任务执行过程中有错误，则
					case p.errChan <- err:
						// if atomic.LoadInt32(&p.closed) != 1 {
						// mylog.Error(err)
						atomic.StoreInt32(&p.closed, 1)
					default:
					}
				}
			}
		}()
	}
}

// templates 使用WorkerPool的例子
func TestWorkerPoolStart(t *testing.T) {
	wp := New(10) // Set the maximum number of threads
	wp.SetTimeout(time.Millisecond)
	for i := 0; i < 20; i++ { // Open 20 requests
		ii := i
		// 添加任务到worker_pool的等待队列里面
		wp.Do(func() error {
			for j := 0; j < 10; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Millisecond)
			}
			// time.Sleep(1 * time.Second)
			return nil
		})
	}
	// 阻塞 等待任务执行完成
	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}
