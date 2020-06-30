package wokerpool

import (
	"fmt"
	"github.com/pbenner/threadpool"
	"math/rand"
	"sync"
	"time"
)

/**
  1、场景描述
     假设有一个任务，分成A、B、C、D四个步骤，四个步骤的耗时差别很大，且不同的任务可能是B的耗时最长，也有可能是D的耗时最长，
     步骤B和C依赖步骤A，步骤D依赖B和C。为了提高性能，故实现任务之间的并发。
2、具体实现
   用四个队列分别完成任务中的每个步骤，队列之间是并发的，队列中可以顺序执行也可以并发执行（比如queue_B）
*/

//
//         b
//    a        d
//         c
//
//
type Product struct {
	name   string
	result string
	isStop bool
}

func queue_A(wg *sync.WaitGroup, A_B chan<- Product, A_C chan<- Product, app []string) {
	defer wg.Done()

	for i, appname := range app {
		product := Product{name: appname, result: "success", isStop: false}
		if i == len(app)-1 {
			product.isStop = true
		}
		A_B <- product
		A_C <- product
		fmt.Println("任务A：", appname)
		time.Sleep(time.Duration(200+rand.Intn(1000)) * time.Millisecond)
	}
}

func queue_B(wg *sync.WaitGroup, A_B <-chan Product, B_D chan<- Product, app []string) {
	defer wg.Done()
	pool := threadpool.New(3, 100)

	// jobs are always grouped, get a new group index
	g := pool.NewJobGroup()
	isStop := false
	for {
		if isStop {
			break
		}
		com := <-A_B // 从a_b管道接受任务进行处理
		isStop = com.isStop

		time.Sleep(time.Duration(200+rand.Intn(1000)) * time.Millisecond)
		pool.AddJob(g, func(pool threadpool.ThreadPool, erf func() error) error {
			i := pool.GetThreadId() + 1
			fmt.Println(i)
			return task(com, B_D)
		})
		// wait until all jobs in group g are done, meanwhile, this thread
		// is also used as a worker

	}
	pool.Wait(g)

}
func task(com Product, B_D chan<- Product) error {
	time.Sleep(3 * time.Second)
	fmt.Println("任务B：", com.name)

	product := Product{name: com.name, result: "success", isStop: false}
	if com.isStop {
		product.isStop = true
	}
	B_D <- com
	return nil
}
func queue_C(wg *sync.WaitGroup, A_C <-chan Product, C_D chan<- Product, app []string) {
	defer wg.Done()
	isStop := false
	for {
		if isStop {
			break
		}
		pvc := <-A_C
		isStop = pvc.isStop
		time.Sleep(time.Duration(200+rand.Intn(1000)) * time.Millisecond)

		fmt.Println("任务C：", pvc.name)

		product := Product{name: pvc.name, result: "success", isStop: false}
		if pvc.isStop {
			product.isStop = true
		}
		C_D <- pvc
	}

}

func queue_D(wg *sync.WaitGroup, C_D <-chan Product, B_D <-chan Product, app []string) {
	defer wg.Done()
	isStop := false
	for {
		if isStop {
			break
		}
		pvc := <-C_D
		com := <-B_D
		if pvc.isStop || com.isStop {
			isStop = true
		}
		time.Sleep(time.Duration(200+rand.Intn(1000)) * time.Millisecond)
		fmt.Println("任务D：", pvc.name, com.name)

	}

}
func main() {
	wgp := &sync.WaitGroup{}
	wgp.Add(4)
	app := []string{"app1", "app2", "app3", "app4", "app5", "app6", "app7", "app8", "app9"}
	A_B := make(chan Product, len(app))
	A_C := make(chan Product, len(app))
	B_D := make(chan Product, len(app))
	C_D := make(chan Product, len(app))
	defer close(A_B)
	defer close(A_C)
	defer close(B_D)
	defer close(C_D)

	go queue_A(wgp, A_B, A_C, app)
	go queue_B(wgp, A_B, B_D, app)
	go queue_C(wgp, A_C, C_D, app)
	go queue_D(wgp, C_D, B_D, app)
	time.Sleep(time.Duration(1) * time.Second)
	wgp.Wait()
}
