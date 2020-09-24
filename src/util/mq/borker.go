package mq

import (
	"errors"
	"sync"
	"time"
)

// Goroutine 和 Channel 是 Go 语言并发编程的两大基石。Goroutine 用于执行并发任务，
// Channel 用于 goroutine 之间的同步、通信。Go提倡使用通信的方法代替共享内存，
// 当一个Goroutine需要和其他Goroutine资源共享时，Channel就会在他们之间架起一座桥梁，
// 并提供确保安全同步的机制。channel本质上其实还是一个队列，遵循FIFO原则。具体规则如下：
// 先从 Channel 读取数据的 Goroutine 会先接收到数据；
// 先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利；
// Go语言中无缓冲的通道（unbuffered channel）是指在接收前没有能力保存任何值的通道。这种类型的通道要求发送 goroutine 和接收 goroutine 同时准备好，才能完成发送和接收操作。
//无缓冲通道的定义方式如下：
//通道实例 := make(chan 通道类型)

type Broker interface {
	// 这些代码我都定义的是内部方法，也就是包外不可用 【函数名小写】
	publish(topic string, msg interface{}) error               // 进行消息的推送，有两个参数即topic、msg，分别是订阅的主题、要传递的消息
	subscribe(topic string) (<-chan interface{}, error)        // 消息的订阅，传入订阅的主题，即可完成订阅，并返回对应的channel通道用来接收数据
	unsubscribe(topic string, sub <-chan interface{}) error    // 取消订阅，传入订阅的主题和对应的通道
	close()                                                    // 这个的作用就是很明显了，就是用来关闭消息队列的
	broadcast(msg interface{}, subscribers []chan interface{}) // 这个属于内部方法，作用是进行广播，对推送的消息进行广播，保证每一个订阅者都可以收到
	setConditions(capacity int)                                // 这里是用来设置条件，条件就是消息队列的容量，这样我们就可以控制消息队列的大小了
}

// broker的实现
type BrokerImpl struct {
	exit     chan bool // 也是一个通道，这个用来做关闭消息队列用的
	capacity int       // 即用来设置消息队列的容量
	// 这里使用一个map结构，key即是topic，其值则是一个切片，chan类型，这里这么做的原因是我们一个topic可以有多个订阅者，所以一个订阅者对应着一个通道
	topics map[string][]chan interface{} // key： topic  value ： queue
	// 读写锁，这里是为了防止并发情况下，数据的推送出现错误，所以采用加锁的方式进行保证
	sync.RWMutex // 同步锁
}

func NewBroker() *BrokerImpl {
	return &BrokerImpl{
		exit:   make(chan bool),
		topics: make(map[string][]chan interface{}),
	}
}

// 一个是设置我们的消息队列容量
func (b *BrokerImpl) setConditions(capacity int) {
	b.capacity = capacity
}

func (b *BrokerImpl) close() {
	select {
	case <-b.exit:
		return
	default:
		close(b.exit)
		b.Lock()
		// 这句代码b.topics = make(map[string][]chan interface{})比较重要，这里主要是为了保证下一次使用该消息队列不发生冲突
		b.topics = make(map[string][]chan interface{})
		b.Unlock()
	}
	return
}

func (b *BrokerImpl) publish(topic string, pub interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return nil
	}

	b.broadcast(pub, subscribers)
	return nil
}

func (b *BrokerImpl) broadcast(msg interface{}, subscribers []chan interface{}) {
	count := len(subscribers)
	concurrency := 1
	// 考虑这样一个问题，当有大量的订阅者时，，比如10000个，我们一个for循环去做消息的推送，
	// 那推送一次就会耗费很多时间，并且不同的消费者之间也会产生延时，，所以采用这种方法进行分解可以降低一定的时间
	switch {
	case count > 1000:
		concurrency = 3
	case count > 100:
		concurrency = 2
	default:
		concurrency = 1
	}

	//采用Timer 而不是使用time.After 原因：time.After会产生内存泄漏 在计时器触发之前，垃圾回收器不会回收Timer
	// 在推送的时候，当推送失败时，我们也不能一直等待呀，所以这里我们加了一个超时机制，超过5毫秒就停止推送，接着进行下面的推送。
	idleDuration := 5 * time.Millisecond
	idleTimeout := time.NewTimer(idleDuration)
	defer idleTimeout.Stop()
	pub := func(start int) {
		for j := start; j < count; j += concurrency {
			idleTimeout.Reset(idleDuration)
			select {
			case subscribers[j] <- msg: // 推送成功
			case <-idleTimeout.C: // 超时
			case <-b.exit: // 退出
				return
			}
		}
	}
	for i := 0; i < concurrency; i++ {
		go pub(i) // 函数是一等公民 。。
	}
}

// 这里的实现则是为订阅的主题创建一个channel，然后将订阅者加入到对应的topic中就可以了，并且返回一个接收channel。
func (b *BrokerImpl) subscribe(topic string) (<-chan interface{}, error) {
	select {
	case <-b.exit:
		return nil, errors.New("broker closed")
	default:
	}

	ch := make(chan interface{}, b.capacity)
	b.Lock()
	b.topics[topic] = append(b.topics[topic], ch)
	b.Unlock()
	return ch, nil
}

// 这里实现的思路就是将我们刚才添加的channel删除就可以了
func (b *BrokerImpl) unsubscribe(topic string, sub <-chan interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()

	if !ok {
		return nil
	}
	// delete subscriber
	b.Lock()
	var newSubs []chan interface{}
	for _, subscriber := range subscribers {
		if subscriber == sub {
			continue
		}
		newSubs = append(newSubs, subscriber)
	}

	b.topics[topic] = newSubs
	b.Unlock()
	return nil
}
