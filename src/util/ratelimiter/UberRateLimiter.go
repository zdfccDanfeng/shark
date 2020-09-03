package ratelimiter

import (
	"fmt"
	"go.uber.org/ratelimit"
	"time"
)

func TestUberRateLimiter() {
	rl := ratelimit.New(100) // per second

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev)) // 耗时
		prev = now
	}

	// Output:
	// 0 0
	// 1 10ms
	// 2 10ms
	// 3 10ms
	// 4 10ms
	// 5 10ms
	// 6 10ms
	// 7 10ms
	// 8 10ms
	// 9 10ms
}

// ticker定时器表示每隔一段时间就执行一次，一般可执行多次。
// timer定时器表示在一段时间后执行，默认情况下只执行一次，如果想再次执行的话，每次都需要调用 time.Reset()方法，
// 此时效果类似ticker定时器。同时也可以调用stop()方法取消定时器
//timer定时器比ticker定时器多一个Reset()方法，两者都有Stop()方法，表示停止定时器,底层都调用了stopTimer()函数。
func TestTicker() {
	// Ticker 包含一个通道字段C，每隔时间段 d 就向该通道发送当时系统时间。
	// 它会调整时间间隔或者丢弃 tick 信息以适应反应慢的接收者。
	// 如果d <= 0会触发panic。关闭该 Ticker 可以释放相关资源。

	ticker1 := time.NewTicker(5 * time.Second)
	// 一定要调用Stop()，回收资源
	defer ticker1.Stop()
	go func(t *time.Ticker) {
		for {
			// 每5秒中从chan t.C 中读取一次
			<-t.C
			fmt.Println("Ticker:", time.Now().Format("2006-01-02 15:04:05"))
		}
	}(ticker1)

	time.Sleep(30 * time.Second)
	fmt.Println("ok")
	// 执行结果
	//
	//开始时间： 2020-03-19 17:49:41
	//Ticker: 2020-03-19 17:49:46
	//Ticker: 2020-03-19 17:49:51
	//Ticker: 2020-03-19 17:49:56
	//Ticker: 2020-03-19 17:50:01
	//Ticker: 2020-03-19 17:50:06
	//结束时间： 2020-03-19 17:50:11
	//ok
}

func TestTimer() {
	// NewTimer 创建一个 Timer，它会在最少过去时间段 d 后到期，向其自身的 C 字段发送当时的时间
	timer1 := time.NewTimer(5 * time.Second)

	fmt.Println("开始时间：", time.Now().Format("2006-01-02 15:04:05"))
	go func(t *time.Timer) {
		times := 0
		for {
			<-t.C
			fmt.Println("timer", time.Now().Format("2006-01-02 15:04:05"))

			// 从t.C中获取数据，此时time.Timer定时器结束。如果想再次调用定时器，只能通过调用 Reset() 函数来执行
			// Reset 使 t 重新开始计时，（本方法返回后再）等待时间段 d 过去后到期。
			// 如果调用时 t 还在等待中会返回真；如果 t已经到期或者被停止了会返回假。
			times++
			// 调用 reset 重发数据到chan C
			fmt.Println("调用 reset 重新设置一次timer定时器，并将时间修改为2秒")
			t.Reset(2 * time.Second)
			if times > 3 {
				fmt.Println("调用 stop 停止定时器")
				t.Stop()
			}
		}
	}(timer1)

	time.Sleep(30 * time.Second)
	fmt.Println("结束时间：", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("ok")
	// 执行结果
	//
	//开始时间： 2020-03-19 17:41:59
	//timer 2020-03-19 17:42:04
	//调用 reset 重新设置一次timer定时器，并将时间修改为2秒
	//timer 2020-03-19 17:42:06
	//调用 reset 重新设置一次timer定时器，并将时间修改为2秒
	//timer 2020-03-19 17:42:08
	//调用 reset 重新设置一次timer定时器，并将时间修改为2秒
	//timer 2020-03-19 17:42:10
	//调用 reset 重新设置一次timer定时器，并将时间修改为2秒
	//调用 stop 停止定时器
	//结束时间： 2020-03-19 17:42:29
	//ok
	//可以看到，第一次执行时间为5秒以后。然后通过调用 time.Reset() 方法再次激活定时器，定时时间为2秒，
	//最后通过调用time.Stop()把前面的定时器取消掉
}
