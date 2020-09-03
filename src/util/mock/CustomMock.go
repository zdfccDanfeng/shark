package mock

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// timers represents a list of sortable timers.
type Timers []*Timer

func (ts Timers) Len() int { return len(ts) }

func (ts Timers) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts Timers) Less(i, j int) bool {
	return ts[i].Next().Before(ts[j].Next())
}

func (ts *Timers) Push(t interface{}) {
	*ts = append(*ts, t.(*Timer))
}

func (ts *Timers) Pop() interface{} {
	t := (*ts)[len(*ts)-1]
	*ts = (*ts)[:len(*ts)-1]
	return t
}

// Mock represents a mock clock that only moves forward programmically.
// It can be preferable to a real-time clock when testing time-based functionality.
type Mock struct {
	sync.Mutex
	now    time.Time // current time
	timers Timers    // timers
}

// NewMock returns an instance of a mock clock.
// The current time of the mock clock on initialization is the Unix epoch.
func NewMock() *Mock {
	return &Mock{now: time.Now()}
}

// Add moves the current time of the mock clock forward by the duration.
// This should only be called from a single goroutine at a time.
func (m *Mock) Add(d time.Duration) {
	m.Lock()
	// Calculate the final time.
	end := m.now.Add(d)

	for len(m.timers) > 0 && m.now.Before(end) {
		t := heap.Pop(&m.timers).(*Timer)
		m.now = t.next
		m.Unlock()
		fmt.Println("=== ^^^^^ +++++")
		t.Tick()
		m.Lock()
	}

	m.Unlock()
	// Give a small buffer to make sure the other goroutines get handled.
	nap()
}

// Timer produces a timer that will emit a time some duration after now.
func (m *Mock) Timer(d time.Duration) *Timer {
	ch := make(chan time.Time)
	fmt.Printf("ch is %v, ch size is : %d\n", ch, len(ch))
	fmt.Printf("m is: now: %v, timers is : %v \n", m.now, m.timers)
	// 初始化的时候 两个使用的是同一个channel 引用。。。
	t := &Timer{
		C:    ch,
		c:    ch,
		mock: m,
		next: m.now.Add(d),
	}
	fmt.Printf("diff is %v\n", t.next.Sub(m.now).Seconds())
	m.addTimer(t)
	fmt.Printf("m is: now: %v, timersSize is : %d \n", m.now, len(m.timers))
	fmt.Printf("detail Timer is : %v,  %v,  %v \n", m.timers[0].mock, m.timers[0].C, m.timers[0].next)
	return t
}

func (m *Mock) addTimer(t *Timer) {
	m.Lock()
	defer m.Unlock()
	heap.Push(&m.timers, t)
}

// After produces a channel that will emit the time after a duration passes.
func (m *Mock) After(d time.Duration) <-chan time.Time {
	return m.Timer(d).C
}

// AfterFunc waits for the duration to elapse and then executes a function.
// A Timer is returned that can be stopped.
func (m *Mock) AfterFunc(d time.Duration, f func()) *Timer {
	t := m.Timer(d) // 定义定时器
	fmt.Println("==== add timer ===")
	fmt.Printf("m now is: %v\n", m.Now())
	go func() {
		fmt.Println("---- 差点错过 ----")
		fmt.Printf("c size is %d\n", len(t.c))
		<-t.c
		f()
	}()
	nap()
	return t
}

// Now returns the current wall time on the mock clock.
func (m *Mock) Now() time.Time {
	m.Lock()
	defer m.Unlock()
	return m.now
}

// Sleep pauses the goroutine for the given duration on the mock clock.
// The clock must be moved forward in a separate goroutine.
// 阻塞等待直到收到消息为止。。。。！！ todo
func (m *Mock) Sleep(d time.Duration) {
	<-m.After(d)
}

// Timer represents a single event.
type Timer struct {
	C    <-chan time.Time
	c    chan time.Time // 接收next信息。。。
	next time.Time      // next tick time
	mock *Mock          // mock clock
}

func (t *Timer) Next() time.Time { return t.next }

func (t *Timer) Tick() {
	select {
	case t.c <- t.next:
		// 收到next元素通知。。。
	default:
	}
	nap()
}

// 暂停 使得其他goroutines有运行的机会
// Sleep momentarily so that other goroutines can process.
func nap() { time.Sleep(1 * time.Millisecond) }
