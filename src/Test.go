package main

import (
	"container/list"
	"context"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/shark/src/config"
	"github.com/shark/src/dao"
	"github.com/shark/src/util"
	"github.com/shark/src/util/lsm"
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
	lsm.TestKeyDb()
}
