package speed

import (
	"fmt"
	"time"
)

type Speed struct {
	interval time.Duration
	ch       chan struct{}
}

func New(interval time.Duration) *Speed {
	s := &Speed{interval: interval,
		ch: make(chan struct{}, 1)}
	go s.timer()
	return s
}

func (s *Speed) timer() {
	for range time.Tick(s.interval) {
		fmt.Println("=========")
		fmt.Printf("chSize is : %d\n", len(s.ch))
		<-s.ch // 出队列。。。
		fmt.Printf("After chSize is : %d\n", len(s.ch))

	}
}

func (s *Speed) Wait() {
	s.ch <- struct{}{}
}
