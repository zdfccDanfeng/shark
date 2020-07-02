package skiplist

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// https://studygolang.com/articles/22491
// Skip List ，称之为跳表，它是一种可以替代平衡树的数据结构，其数据元素默认按照key值升序，天然有序。
//Skip list让已排序的数据分布在多层链表中，以0-1随机数决定一个数据的向上攀升与否，通过“空间来换取时间”的一个算法，
//在每个节点中增加了向前的指针，在插入、删除、查找时可以忽略一些不可能涉及到的结点，从而提高了效率。

// =======================================================================================================
// SkipList具备如下特性：
//
//由很多层结构组成，level是通过一定的概率随机产生的
//每一层都是一个有序的链表，默认是升序，也可以根据创建映射时所提供的Comparator进行排序，具体取决于使用的构造方法
//最底层(Level 1)的链表包含所有元素
//如果一个元素出现在Level i 的链表中，则它在Level i 之下的链表也都会出现
//每个节点包含两个指针，一个指向同一链表中的下一个元素，一个指向下面一层的元素
//1、层数（LEVEL）越高，链上的节点越少，大致呈P=0.5的几何分布。
// 2、每一层的节点均有序且不重复。
//

const (
	// DefaultMaxLevel 默认skip list最大深度
	DefaultMaxLevel int = 18
	// DefaultProbability 默认的概率
	DefaultProbability float64 = 1 / math.E
)

// elementNode 数组指针，指向元素
type elementNode struct {
	next []*Element
}

// Element 跳转表数据结构 skipList里面的一个元素结构
type Element struct {
	elementNode
	key   float64     // 用以排序和判断大小的关键字
	value interface{} // 定义元素 附加
}

// Key 获取key的值
func (e *Element) Key() float64 {
	return e.key
}

// Value 获取key的值
func (e *Element) Value() interface{} {
	return e.value
}

// 定义整体的SkipList
type SkipList struct {
	elementNode              // 指针数组
	maxLevel    int          // 最大深度
	length      int          // 长度
	randSource  rand.Source  // 随机数种子，动态调节跳转表的长度
	probability float64      // 概率 来初始化新添加进来的元素的 level, 每个元素都有个一个level, 也就是层数从高到低，数量按比例增多
	probTable   []float64    // 存储位置，对应key
	mutex       sync.RWMutex // 保证线程安全
	// prevNodesCache 的作用在哪里，这就是用来插入新元素的时候把元素中的指针数组把前后相同层的元素连接起来
	prevNodesCache []*elementNode // 缓存
}

// NewSkipList 新建跳转表
func NewSkipList() *SkipList {
	return NewWithMaxLevel(DefaultMaxLevel)
}

// ProbabilityTable 初始化 Probability Table
func ProbabilityTable(probability float64, maxLevel int) (table []float64) {
	for i := 1; i <= maxLevel; i++ {
		prob := math.Pow(probability, float64(i-1))
		table = append(table, prob)
	}
	return table
}

// NewWithMaxLevel 自定义maxLevel新建跳转表
func NewWithMaxLevel(maxLevel int) *SkipList {
	if maxLevel < 1 || maxLevel > DefaultMaxLevel {
		panic("invalid maxlevel")
	}

	return &SkipList{
		elementNode:    elementNode{next: make([]*Element, maxLevel)}, // 需要初始化指针数组，指针数组的长度就是跳表的层次
		prevNodesCache: make([]*elementNode, maxLevel),
		maxLevel:       maxLevel,
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
		probability:    DefaultProbability,
		probTable:      ProbabilityTable(DefaultProbability, maxLevel),
	}
}

// 随机计算最接近的
// 这个主要是插入一个元素时，我们要给这个元素指定高度，那就是通过这个函数在概率的范围下指定这个元素的高度。后面 添加元素的函数会再次提现。
func (list *SkipList) randLevel() (level int) {
	r := float64(list.randSource.Int63()) / (1 << 63)
	level = 1
	for level < list.maxLevel && r < list.probTable[level] {
		level++ // 级别追加
	}

	return level
}

// SetProbability 设置新的概率,刷新概率表
func (list *SkipList) SetProbability(newProbability float64) {
	list.probability = newProbability
	list.probTable = ProbabilityTable(newProbability, list.maxLevel)
}

// Set 存储新的值
func (list *SkipList) Set(key float64, value interface{}) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock() // 线程安全

	var element *Element
	prevs := list.getPrevElementNodes(key)
	if element = prevs[0].next[0]; element != nil && key == element.key {
		element.value = value
		return element
	}

	element = &Element{
		elementNode: elementNode{next: make([]*Element, list.randLevel())},
		key:         key,
		value:       value,
	}
	list.length++

	for i := range element.next { // 插入数据
		element.next[i] = prevs[i].next[i]
		prevs[i].next[i] = element // 记录位置
	}

	return element
}

// Get 获取key对应的值
// 在图中的第一个查找中，查询12
//
//首先会从最左边最上层开始向右遍查找，然后会找到 6
//然后 根据 6 的再下一层 也就是第三层向后找，没有找到
//然后 继续根据 6 的第二层开始向后找，找到了 9
//继续在 9 的更下一层继续找，然后可以找到 12
// 总体思想：从左到右，从上到下。
func (list *SkipList) Get(key float64) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock() // 线程安全

	var prev *elementNode = &list.elementNode // 保存前置结点 ，因为如果续后没有找到合适的，需要通过前置找下一层 down !!
	var next *Element
	// 从最顶层开始依次往下面进行探索寻找
	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i] // 循环跳到下一个
		for next != nil && key > next.key {
			prev = &next.elementNode // 继续向右边探索。
			next = next.next[i]
		}
	}

	if next != nil && next.key == key { // 找到
		return next
	}

	return nil // 没有找到
}

// Remove 获取key对应的值
func (list *SkipList) Remove(key float64) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock() // 线程安全
	// 最核心代码位置：
	var element *Element
	prevs := list.getPrevElementNodes(key)
	// 在 Set 函数中拿到 prevs 这样的一个数组后，在查找元素12中就可以通过prevs[0].next[0]可以索引到元素12，判断是否找到元素
	if element = prevs[0].next[0]; element != nil && key == element.key {
		for k, v := range element.next {
			prevs[k].next[k] = v // 删除
		}

		list.length--
		return element
	}

	return nil
}

// 用来记录我们在查找 key 的中途会经过的元素的指针数组
// 在上面的查询元素12的例子中，prevs 记录的值分别就是 :
//
// prevs[3] == 6.elementCode
// prevs[2] == 6.elementCode
// prevs[1] == 9.elecmentCode
// prevs[0] == 9.elementCode
func (list *SkipList) getPrevElementNodes(key float64) []*elementNode {
	var prev *elementNode = &list.elementNode // 保存前置结点
	var next *Element
	prevs := list.prevNodesCache // 缓冲集合
	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i] // 循环跳到下一个
		for next != nil && key > next.key {
			prev = &next.elementNode // 继续向右边找。。
			next = next.next[i]
		}
		prevs[i] = prev
	}
	return prevs
}
