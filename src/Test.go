package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/scaffold/src/config"
	"github.com/scaffold/src/dao"
	"github.com/scaffold/src/util"
	"github.com/scaffold/src/util/algorithm/dfs"
	"sort"
	"strconv"
	"strings"
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

	for index, _ := range nums {
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
	dfs.TestLongestLenght(")(")
}
