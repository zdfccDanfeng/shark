package customworkerpool

import (
	"fmt"
	"github.com/shark/src/util/customworkerpool/workerpool"
	"runtime"
	"strconv"
	"time"
)

// 我创建了一个 MyTask 的类型，它定义了工作执行的状态。接着我实现一个 MyTask 的函数成员 DoWork，
// 它同时符合 PoolWorker 接口的签名。由于 MyTask 实现了 PoolWorker 的接口，MyTask 类型的对象也被认为是 PoolWorker 类型的对象。
// 现在我们把 MyTask 类型的对象传入 PostWork 方法中。
// 实现PoolWorker

// 我设置 Go 运行环境使用我本机上的全部 CPU 和核心，我创建了一个 24 个 Goroutines 的工作池。我本机有 8 个核心，
//就像上面我们得到的结论，每个核心分配 3 个 Goroutines 是比较好的配置。最后一个参数是告诉工作池创建一个容量为 100 个任务的队列。
//
// 接着我创建了一个 MyTask 的对象并且提交到队列中。为了记录日志，PostWork 方法的第一个参数可以设置成调用方的名称。
// 如果调用返回的 err 参数是空，表明此任务已经得到提交；如果非空，大概率意味着已经超过了队列的容量，你的任务未能得到提交。
type MyTask struct {
	Name string
	WP   *workerpool.WorkPool // 绑定一个任务池，任务的执行交给任务池负责
}

func (mt *MyTask) DoWork(workRoutineId int) {
	fmt.Println(mt.Name)

	fmt.Printf("*******> WR: %d QW: %d AR: %d\n",
		workRoutineId,
		mt.WP.QueuedWork(),
		mt.WP.ActiveRoutines())

	time.Sleep(100 * time.Millisecond)
}

func TestRun(i int) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	workPool := workerpool.New(runtime.NumCPU()*3, 100)

	task := MyTask{
		Name: "A" + strconv.Itoa(i),
		WP:   workPool,
	}

	err := workPool.PostWork("main", &task)
	if err != nil {
		panic(err)
	}
	// …
}
