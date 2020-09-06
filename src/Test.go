package main

import (
	"container/list"
	"context"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/shark/src/config"
	"github.com/shark/src/dao"
	"github.com/shark/src/rpc"
	"github.com/shark/src/util"
	"github.com/shark/src/util/algorithm/tree"
	"github.com/shark/src/util/mock"
	"github.com/shark/src/util/speed"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func largestPerimeter(A []int) int {
	len := len(A)
	if len <= 2 {
		return 0
	}
	res := 0
	sort.Ints(A)
	for index := len - 1; index > 1; index-- {
		a, b, c := A[index], A[index-1], A[index-2]
		if b+c > a {
			res = max(a+b+c, res)
		} else {
			// a 太大了
			continue
		}

	}
	return res
}

func canJump(nums []int) bool {

	if len(nums) <= 1 {
		return true
	}

	reach := 0 // 可以抵达的位置，每一步尽可能的走远

	length := len(nums)

	for index := range nums {
		if index > reach || reach >= length-1 {
			// index > reach的意义：表征前面的累积reach 努力无法抵达当前的位置，因此直接退出这种选择
			break
		}
		reach = max(reach, index+nums[index]) // 取最远抵达
	}

	return reach >= length-1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type IndustryMapping struct {
	Id           int64  `json:"id" form:"id"`
	ParentId     int64  `json:"parent_id" form:"parent_id"`
	IndustryName string `json:"industry_name" form:"industry_name"`
	Level        int    `json:"level" form:"level"`
}

//指定IndustryMapping结构体对应的数据表为ad_dsp_industry_v4
func (p IndustryMapping) TableName() string {
	return "ad_dsp_industry_v4"
}

//  // Len is the number of elements in the collection.
//    Len() int
//    // Less reports whether the element with
//    // index i should sort before the element with index j.
//    Less(i, j int) bool
//    // Swap swaps the elements with indexes i and j.
//    Swap(i, j int)

type Industry_Mappins []IndustryMapping

func (p Industry_Mappins) Len() int {
	return len(p)
}
func (p Industry_Mappins) Less(i, j int) bool {
	return p[i].Id < p[j].Id
}
func (p Industry_Mappins) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func ReadFiles(path string) *Industry_Mappins {
	data := util.TestReadBigFile(path)
	mappings := make([]IndustryMapping, 0)
	set := hashset.New()
	for index, line := range data {
		if len(line) == 0 || index == 0 {
			continue
		}
		data := strings.Split(line, ",")
		second_industry_id, _ := strconv.ParseInt(data[0], 10, 64)
		second_industry_name := data[1]
		first_industry_name := data[2]
		first_indusry_id, _ := strconv.ParseInt(data[3], 10, 64)
		if !set.Contains(second_industry_id) {
			mapping := IndustryMapping{Id: second_industry_id, ParentId: first_indusry_id, IndustryName: second_industry_name, Level: 2}
			mappings = append(mappings, mapping)
			set.Add(second_industry_id)
			if !set.Contains(first_indusry_id) {
				industryMapping := IndustryMapping{Id: first_indusry_id, ParentId: -1, IndustryName: first_industry_name, Level: 1}
				mappings = append(mappings, industryMapping)
				set.Add(first_indusry_id)
			}
		}
	}
	res := Industry_Mappins(mappings)
	sort.Sort(res)
	return &res
}

func InsertIndustryMapings(p *Industry_Mappins) {
	conf := config.Config().Dbs["online"]
	conn := dao.InitConn(conf.Username, conf.Password, conf.Host, conf.Database, conf.Port)
	for _, data := range *p {
		conn.Create(data)
	}
}

/**
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `age` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '年龄',
  `gender` varchar(5) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `platform` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '系统版本',
  `region` mediumtext COLLATE utf8mb4_unicode_ci COMMENT '区域定向',
  `region_ids` mediumtext COLLATE utf8mb4_unicode_ci COMMENT '地域ids',
  `language` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '语言',
  `fans_star` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '网红粉丝',
  `interest_video` varchar(1023) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '兴趣视频',
  `business_interest` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '商业兴趣',
  `business_interest_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '商业兴趣类型，0-不限，1-智能定向，2-兴趣标签',
  `device_price` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设备价格',
  `device_brand` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设备品牌',
  `package_name` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'packageName定向',
  `page` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `network` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `interest` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `audience` varchar(511) COLLATE utf8mb4_unicode_ci DEFAULT '[]' COMMENT '人群定向包, [1, 2, 3]',
  `paid_audience` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT '[]' COMMENT '付费人群包, [1, 2, 3]',
  `md5` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `population` varchar(3000) COLLATE utf8mb4_unicode_ci DEFAULT '[]' COMMENT '广告主上传的人群定向包,格式为ID的集合，如[1,2,3]',
  `exclude_population` varchar(2000) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '[]' COMMENT '定向排除人群包',
  `intelli_extend` varchar(1023) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '{}' COMMENT '智能扩量',
  `social_star_label` text COLLATE utf8mb4_unicode_ci COMMENT '网红标签',
  `behavior_interest_keyword` text COLLATE utf8mb4_unicode_ci COMMENT '行为兴趣标签',

*/
type AdDspTarget struct {
	id       int64
	age      string
	gender   string
	platform string
	region   string
}

func TestDuopleList() {
	doubleList := list.New()
	doubleList.PushFront(10) // 头插法
	doubleList.PushFront(20)
	doubleList.PushFront(45)
	doubleList.PushFront(67)
	font := doubleList.Front()
	fmt.Println("font is :", font)
	back := doubleList.Back()
	doubleList.MoveToFront(back)
}

// 协程的生命周期：创建 、回收（系统、gc自动完成）、中断（主要通过context来实现）
func TestContext() {
	parent := context.Background() // 初始化context
	// 生成一个取消的context
	ctx, cancel := context.WithCancel(parent)
	runTimes := 0
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("GoRoutine is done!!!")
				return
			default:
				fmt.Printf("GoRoutine Running times : %d\n", runTimes)
				runTimes += 1
			}
			if runTimes >= 5 {
				cancel() // 关闭GoRoutine
				wg.Done()
			}
		}
	}(ctx)
	wg.Wait()
}

// 100 人抢10个鸡蛋 ,利用channel原理实现资源争抢
func GetEggs() {
	var wg sync.WaitGroup
	eggs := make(chan int, 10)
	for i := 0; i < 10; i++ {
		eggs <- i
	}
	for i := 0; i < 100; i++ {
		go func(num int) {
			wg.Add(1)
			select {
			case egg := <-eggs:
				{
					fmt.Printf("people %d get egg %d \n", num, egg)
				}
			default:
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

var (
	infos = make(chan int, 10)
	// 声明全局变量，节省内存交互时间
	wg     sync.WaitGroup
	global sync.WaitGroup
	ch     = make(chan []int, 1)
)

// A 车，清洗材料
func funcA(elemetns []int) {
	var tasks = make([][]int, 3)
	// 3个工人
	for i := 0; i < 3; i++ {
		// 用来分配任务
		task := []int{}
		// 获取任务分割
		for _, value := range elemetns {
			task = append(task, value/3.0)
		}
		wg.Add(1)
		go func(task []int) {
			tasks[i] = clean(task)
			wg.Done()
		}(task)
		wg.Wait()

		// 合并回elements
		for index, _ := range tasks {
			// 清空原来的elements
			elemetns[index] = 0
			for _, task := range tasks {
				elemetns[index] += task[index]
			}
		}
		ch <- elemetns
		global.Done()
	}
}

func MainPipline() {
	elemetns := []int{1, 1, 1}
	global.Add(1)
	go funcA(elemetns)
	global.Add(1)
	go funcB()
	global.Add(1)
	go funcC()
	global.Wait()
}

// 清洗材料
func clean(task []int) []int {
	fmt.Printf("clean task %v\n", task)
	return task
}

// B车， 加工材料
func funcB() {
	elements := []int{}
	for {
		select {
		case elements = <-ch:
			// 阻塞直到接收到材料
			break
		default:
			continue
		}
	}
	for index, element := range elements {
		wg.Add(1)
		go func(element, index int) {
			// 加工材料
			elements[index] = cure(element)
			wg.Done()
		}(element, index)
	}
	wg.Wait()
	global.Done()
}

func cure(element int) int {
	fmt.Printf("cue ele %d\n", element)
	return element
}

// C车，运输材料
func funcC() {

	elements := []int{}
	for {
		select {
		// 对B车的结果进行接收
		case elements = <-ch:
			break
		default:
			continue
		}
	}
	for index, element := range elements {
		wg.Add(1)
		go func(index, element int) {
			elements[index] = carry(element)
		}(element, index)
	}
	wg.Wait()
	global.Done()
}

func carry(element int) int {
	fmt.Printf("carry ele %d\n", element)
	return element
}
func producer(index int) {
	infos <- index
	fmt.Printf("Producer %d, sent %d\n", index, index)
}
func consumer(index int) {
	fmt.Printf("Consumer %d , reveied  msg %d\n", index, <-infos)
}

func TestConsumerAndProducer() {
	// 十个生产者
	for i := 0; i < 10; i++ {
		go producer(i)
	}
	// 十个消费者
	for i := 0; i < 100; i++ {
		go consumer(i)
	}

	time.Sleep(20 * time.Second)
}

type Animial interface {
	walk()
}

type Person interface {
	addDock()
	Animial
}

// // Context 提供跨越API的截止时间获取，取消信号，以及请求范围值的功能。
//// 它的这些方案在多个 goroutine 中使用是安全的
//type Context interface {
//    // 如果设置了截止时间，这个方法ok会是true，并返回设置的截止时间
// Deadline() (deadline time.Time, ok bool)
//    // 如果 Context 超时或者主动取消返回一个关闭的channel，如果返回的是nil，表示这个
//    // context 永远不会关闭，比如：Background()
// Done() <-chan struct{}
//    // 返回发生的错误
// Err() error
//    // 它的作用就是传值
// Value(key interface{}) interface{}
//}
func TestContext2() {
	req, _ := http.NewRequest("GET", "https://api.github.com/users/helei112g", nil)
	// 这里设置了超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("request Err", err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	// 上面这段程序就是请求 github 获取用户信息的接口，通过 context 包设置了请求超时时间是 1ms （肯定无法访问到）
	// 执行时我们看到控制台做如下输出：
	//2020/xx/xx xx:xx:xx request Err Get https://api.github.com/users/helei112g: context deadline exceeded
	//exit status 1
}

// 模拟获取订单服务。大概意思是说，有一个获取订单详情的请求，会单独起一个 goroutine 去处理该请求。在该请求内部又有三个分支 goroutine 分别处理订单详情、推荐商品、物流信息；每个分支可能又需要单独调用DB、Redis等存储组件。那么面对这个场景我们需要哪些额外的事情呢？
//
//三个分支 goroutine 可能是对应的三个不同服务，我们想要携带一些基础信息过去，比如：LogID、UserID、IP等；
//每个分支我们需要设置过期时间，如果某个超时不影响整个流程；
//如果主 goroutine 发生错误，取消了请求，对应的三个分支应该也都取消，避免资源浪费；
//简单归纳就是传值、同步信号（取消、超时）。
//由于服务内部不方便模拟，我们简化成函数调用，假设图中所有的逻辑都可以并发调用。现在我们的要求是：
//
//整个函数的超时时间为1s；
//需要从最外层传递 LogID/UserID/IP 信息到其它函数；
//获取订单接口超时为 500ms，由于 DB/Redis 是其内部支持的，这里不进行模拟；
//获取推荐超时是 400ms；
//获取物流超时是 700ms。
//为了清晰，我这里所有接口都返回一个字符串，实际中会根据需要返回不同的结果；请求参数也都只使用了 context。代码如下：
type key int

const (
	userIP = iota
	userID
	logID
)

// timeout: 1s
type Result struct {
	order     string
	logistics string
	recommend string
}

func TestGetOrderInfo() (result *Result, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel() // 设置
	// 从最外层传递 LogID/UserID/IP 信息到其它函数；
	ctx = context.WithValue(ctx, userIP, "127.0.0.1") // 子流程分别持有父Context引用！！！！
	ctx = context.WithValue(ctx, userID, 666888)
	ctx = context.WithValue(ctx, logID, "123456")
	result = &Result{} // 业务逻辑处理放到协程
	go func() { result.order, err = getOrderDetail(ctx) }()
	go func() { result.logistics, err = getLogisticsDetail(ctx) }()
	go func() { result.recommend, err = getRecommend(ctx) }()
	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err() // 取消或者超时，把现有已经拿到的结果返回
		default:
		}
		//有错误直接返回
		if err != nil {
			return result, err
		}
		// 全部处理完成，直接返回
		if result.order != "" && result.logistics != "" && result.recommend != "" {
			return result, nil
		}
	}
}

// 获取订单接口超时为 500ms
func getOrderDetail(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()
	// 模拟超时
	time.Sleep(time.Millisecond * 700)
	// 获取 user id
	uip := ctx.Value(userIP).(string)
	fmt.Println("userIP", uip)
	return handleTimeout(ctx, func() string {
		return "order"
	})
}

// 获取物流超时是 700ms。
func getLogisticsDetail(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*700)
	defer cancel()
	// 获取 user id
	uid := ctx.Value(userID).(int)
	fmt.Println("userID", uid)
	return handleTimeout(ctx, func() string {
		return "logistics"
	})
}

// 获取推荐超时是 400ms；
func getRecommend(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*400)
	defer cancel() // 获取 log id
	lid := ctx.Value(logID).(string)
	fmt.Println("logID", lid)
	return handleTimeout(ctx, func() string {
		return "recommend"
	})
}

// 超时的统一处理代码
func handleTimeout(ctx context.Context, f func() string) (string, error) {
	// 请求之前先去检查下是否超时
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	str := make(chan string)
	go func() {
		// 业务逻辑
		str <- f()
	}()
	// 类似异步回调
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	// str 这个channel接收到了回调信号会返回
	case ret := <-str:
		return ret, nil
	}
}

// ["NumArray","update","sumRange","sumRange","update","sumRange"]
//[[[9,-8]],[0,3],[1,1],[0,1],[1,-3],[0,1]]
// ["NumArray","sumRange","sumRange","sumRange","update","update","update","sumRange","update","sumRange","update"]
//[[[0,9,5,7,3]],[4,4],[2,4],[3,3],[4,5],[1,7],[0,8],[1,2],[1,9],[4,4],[3,4]]
// ["NumArray","sumRange","update","sumRange","sumRange","update","update","sumRange","sumRange","update","update"]
//[[[-28,-39,53,65,11,-56,-65,-39,-43,97]],[5,6],[9,27],[2,3],[6,7],[1,-82],[3,-72],[3,7],[1,8],[5,13],[4,-67]]
func TestSegment() {
	numArray := tree.Constructor([]int{-28, -39, 53, 65, 11, -56, -65, -39, -43, 97})
	fmt.Println("nn :", numArray.SumRange(5, 6))
	numArray.Update(0, 3)
	sumRange := numArray.SumRange(1, 1)
	sumRange1 := numArray.SumRange(0, 1)
	numArray.Update(1, -3)

	su := numArray.SumRange(0, 1)
	fmt.Println("sumRange :", sumRange)
	fmt.Println("sumRange1 :", sumRange1)
	fmt.Println("su :", su)

}

func TestMock() {
	mock := mock.NewMock()
	c := mock.AfterFunc(time.Second*3, func() {
		fmt.Println("hhhhhhh")
	})
	c.Tick()
	fmt.Printf("cSize is %d\n", len(c.C))
	//<-c.C

	time.Sleep(time.Second * 10)
}

// Go之函数直接实现接口
//1.定义一个接口
type Run interface {
	Runing()
}

//2.定义一个函数类型
type Runer func()

//3.让函数直接实现接口
func (self Runer) Runing() {
	self()
}

//调用
var run Runer = Runer(func() {
	fmt.Println("i am runing")
})

// run.Runing()

func TestSpeed() {
	sp := speed.New(time.Second * 1)
	sp.Wait()
	time.Sleep(time.Second * 10)
}

func main() {
	//
	//nums := []int{0, 12, 1, 0, 4}
	//res := canJump(nums)
	//print(res)
	//
	//arr := []int{3, 6, 3, 2}
	//perimeter := largestPerimeter(arr)
	//fmt.Println("res is :", perimeter)
	//fmt.Println("============================================")
	///**
	//
	//	     5、、、、、、、
	//	                   \
	//	     /3 ----------- 1
	//       /                '
	//	  4                 '
	//	   \                '
	//	    2 --------------'
	//*/
	//scheduler.Test_schedule()
	//basePath := util.RelativePath()
	//res := ReadFiles(basePath + "/files/付费人群界面排序.csv")
	//log.Println("data is :", res)
	//InsertIndustryMapings(res)
	//   1
	// /   \
	//2     3
	// \
	//  5
	//node1 := tree.TreeNode{Val: 1}
	//node2 := tree.TreeNode{Val: 2}
	//node3 := tree.TreeNode{Val: 3}
	//node4 := tree.TreeNode{Val: 5}
	//
	//node1.Left = &node2
	//node1.Right = &node3
	//node2.Right = &node4
	//
	//tree.TestBinaryTreePaths(&node1)
	//dfs.TestLongestLenght(")(")
	//TestDuopleList()
	//TestContext()
	//GetEggs()
	//TestConsumerAndProducer()
	//MainPipline()
	//lsm.TestKeyDb()
	//TestSegment()
	//reflect.TestParse()
	//TestMock()
	//TestSpeed()
	//var a int32
	//atomic.AddInt32(&a, 1)
	//fmt.Println("a is : ", a)
	//rpc.NewProductServer()
	rpc.NewProductClient()
}
