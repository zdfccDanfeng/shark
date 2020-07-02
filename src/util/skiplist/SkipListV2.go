package skiplist

import (
	"math/rand"
	"sync"
	"time"
)

type SkipListNode struct {
	key  int
	data interface{}
	next []*SkipListNode // 指针切片 ，比如版本2对应的图中 Head指针对应的切片为：[5, 3, 1, 1]
}

// 定义跳跃表结构体
type SkipListV2 struct {
	head   *SkipListNode // 头节点
	tail   *SkipListNode // 尾节点
	length int           // 数据总量
	level  int           // 层数
	mut    *sync.RWMutex // 互斥锁
	rand   *rand.Rand    // 随机数生成器， 用于生成随机层数
}

// 随机数生成层级
func (list *SkipListV2) randomLevel() int {
	level := 1
	for ; level < list.level && list.rand.Uint32()&0x1 == 1; level++ {
	}
	return level
}

func NewSkipList2(level int) *SkipListV2 {

	list := &SkipListV2{}
	if level <= 0 {
		level = 32
	}
	list.level = level
	list.head = &SkipListNode{next: make([]*SkipListNode, level, level)}
	list.tail = &SkipListNode{}
	list.mut = &sync.RWMutex{}
	list.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	for index := range list.head.next {
		list.head.next[index] = list.tail // 初始化各个层为空链表
	}
	return list
}

// 实现数据插入 3
func (list *SkipListV2) Add(key int, data interface{}) {
	list.mut.Lock()
	defer list.mut.Unlock()
	// 确定插入的深度
	level := list.randomLevel() // 假设得到的level = 3
	// 查找插入的位置
	update := make([]*SkipListNode, level, level)
	node := list.head
	// 沿着level3的链条，开始找第一个比 3大的节点或者 tail节点，记录下该节点的前一个节点和level值
	// 思想：从最高层节点，一步一步的向右边进行比较着走，直到right为null或者右边部分的节点Node的key大于当前key为止，
	// 然后继续往下面走，依次重复该过程，直到down为null为止，即找到了前辈，
	for index := level - 1; index >= 0; index-- {
		// 内循环
		for {
			node1 := node.next[index]
			if node1 == list.tail || node1.key > key {
				update[index] = node1 // 找到第一个插入位置
				break
			} else if node1.key == key {
				// update
				node1.data = data
				return
			} else {
				node = node1 // 继续进化 向右边走。。。
			}
		}
	}
	newNode := &SkipListNode{key: key, data: data, next: make([]*SkipListNode, level, level)}
	for index, node := range update {
		node.next[index], newNode.next[index] = newNode, node.next[index]
	}
	list.length++
}

func (list *SkipListV2) Remove(key int) bool {
	list.mut.Lock()
	defer list.mut.Unlock()
	// 查找删除节点
	node := list.head
	remove := make([]*SkipListNode, list.level, list.level)
	var target *SkipListNode
	for index := len(node.next) - 1; index >= 0; index++ {
		for {

			node1 := node.next[index]
			if node1 == list.tail || node1.key > key {
				break
			} else if node1.key == key {
				remove[index] = node
				target = node1
				break
			} else {
				node = node1
			}
		}
	}
	// 执行删除操作
	if target != nil {
		for index, node1 := range remove {
			if node1 != nil {
				node1.next[index] = target.next[index]
			}
		}
		list.length--
		return true
	}
	return false
}

// 查找
func (list *SkipListV2) Find(key int) interface{} {
	list.mut.RUnlock()
	defer list.mut.RUnlock()
	node := list.head
	for index := len(node.next) - 1; index >= 0; index-- {
		for {
			node1 := node.next[index]
			if node1 == list.tail || node1.key > key {
				break // 只能往下一层走了
			} else if node1.key == key {
				return node1.data
			} else {
				node = node1
			}

		}
	}
	return nil
}

// 获取数据总量的方法
func (list *SkipListV2) Length() int {
	// 允许多个线程同时读，但只要有一个线程在写，其他线程就必须等待：
	list.mut.RUnlock()
	defer list.mut.RUnlock()
	return list.length
}
