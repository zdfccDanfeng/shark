package goroutinePool

import "log"

// 任务的属性应该是一个业务函数
type Task struct {
	f func() error // 函数名f,无参数，返回值error
}

func NewTask(arg_f func() error) *Task {
	return &Task{f: arg_f}
}

func (task *Task) Execute() {
	_ = task.f() // 执行任务中已经绑定好的业务方法
}

type Pool struct {
	EntryChannel chan *Task // 对外的ask入口
	JobsChannel  chan *Task // 内部的Task队列
	workerNum    int        // 携程池里面最大的worker数量
}

// 创建Pool
func NewPool(cap int) *Pool {
	pool := Pool{EntryChannel: make(chan *Task),
		JobsChannel: make(chan *Task),
		workerNum:   cap}
	return &pool
}

// pool绑定干活的方法
func (pool *Pool) worker(workID int) {
	// worker工作，永久从JobsChannel取任务，然后执行任务
	for task := range pool.JobsChannel {
		task.Execute()
		log.Printf("worker Id %d has executed!\n", workID)
	}
}

// Pool绑定携程池工作方法
func (pool *Pool) run() {
	for i := 0; i < pool.workerNum; i++ {
		go pool.worker(i)
	}
	// 从EntryChannel 取任务放进JobsChannel
	for task := range pool.EntryChannel {
		pool.JobsChannel <- task // 添加task优先级排序逻辑
	}
}
