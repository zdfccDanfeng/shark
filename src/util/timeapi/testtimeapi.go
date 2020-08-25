package timeapi

import (
	"fmt"
	"time"
)

// https://aimuke.github.io/go/2019/12/12/go-timer-ticker/
// 每两秒给你的女票发送一个"I Love You!"
func Love() {
	timer := time.NewTimer(2 * time.Second) // 新建一个Timer ，timer定时器，是到固定时间后会执行一次

	for {
		select {
		case <-timer.C:
			fmt.Println("I Love You!")
			timer.Reset(2 * time.Second) // 上一个when执行完毕重新设置
		}
	}
	return
}

// 基于Ticker实现的每隔两秒给你的女票发送一个"I Love You!"
func Love2() {
	//定义一个ticker ， ticker只要定义完成，从此刻开始计时，不需要任何其他的操作，每隔固定时间都会触发
	ticker := time.NewTicker(time.Millisecond * 500)
	//Ticker触发
	go func() {
		for t := range ticker.C {
			fmt.Println(t)
			fmt.Println("I Love You!")
		}
	}()

	time.Sleep(time.Second * 18)
	//停止ticker
	ticker.Stop()
}
