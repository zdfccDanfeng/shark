package slidingwindow

import "math"

// 滑动窗口的最大值
// 给定一个数组 nums 和滑动窗口的大小 k，请找出所有滑动窗口里的最大值。
//
//示例:
//
//输入: nums = [1,3,-1,-3,5,3,6,7], 和 k = 3
//输出: [3,3,5,5,6,7]
//解释:
//
//  滑动窗口的位置                最大值
//---------------               -----
//[1  3  -1] -3  5  3  6  7       3
// 1 [3  -1  -3] 5  3  6  7       3
// 1  3 [-1  -3  5] 3  6  7       5
// 1  3  -1 [-3  5  3] 6  7       5
// 1  3  -1  -3 [5  3  6] 7       6
// 1  3  -1  -3  5 [3  6  7]      7
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/hua-dong-chuang-kou-de-zui-da-zhi-lcof
//著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
// todo 思路：
//     滑动窗口，是因为，遍历的时候，两个指针一前一后夹着的子串（子数组）类似一个窗口，这个窗口大小和范围会随着前后指针的移动发生变化
//     保证右指针每次往前移动一格，每次移动都会有新的一个元素进入窗口，这时条件可能就会发生改变，然后根据当前条件来决定左指针是否移动，
//     以及移动多少格
func maxSlidingWindow(nums []int, k int) []int {
	res := make([]int, 0)
	if nums == nil || len(nums) == 0 {
		// 返回空的切片。。。
		return res
	}
	index := make([]int, 0)
	l, r, n := 0, 0, len(nums)
	currMaxIndex := l
	for ; r < n; r++ {
		if nums[currMaxIndex] < nums[r] {
			currMaxIndex = r
		}
		if r >= k-1 {
			index = append(index, currMaxIndex)
		}
		if (r - l) >= k {
			l += 1
			if currMaxIndex == l {
				// findMax
			}
		}
	}
	return res
}

// todo 滑动窗口算法骨架思想：/* 滑动窗口算法框架 */
// 题目问法大致有这几种：
//    给两个字符串，一长一短，问其中短的是否在长的中满足一定的条件存在，例如：
//    1. 求长的的最短子串，该子串必须涵盖短的的所有字符
//    2. 短的的 anagram 在长的中出现的所有位置
//    给一个字符串或者数组，问这个字符串的子串或者子数组是否满足一定的条件，例如：
//     1.含有少于 k 个不同字符的最长子串
//     2.所有字符都只出现一次的最长子串
// 和双指针题目类似，更像双指针的升级版，滑动窗口核心点是维护一个窗口集，根据窗口集来进行处理
//  核心步骤
//   1。 right 右移
//   2。 收缩
//   3。 left 右移
//   4。 求结果
// 需要变化的地方
//  1、右指针右移之后窗口数据更新
//  2、判断窗口是否要收缩
//  3、左指针右移之后窗口数据更新
//  4、根据题意计算结果
// void slidingWindow(string s, string t) {
//    unordered_map<char, int> need, window;
//    for (char c : t) need[c]++;
//    int left = 0, right = 0;
//    int valid = 0;
//    while (right < s.size()) {
//        // c 是将移入窗口的字符
//        char c = s[right];
//        // 右移窗口
//        right++;
//        // 进行窗口内数据的一系列更新
//        ...
//
//        /*** debug 输出的位置 ***/
//        printf("window: [%d, %d)\n", left, right);
//        /********************/
//
//        // 判断左侧窗口是否要收缩
//        while (window needs shrink) {
//            // d 是将移出窗口的字符
//            char d = s[left];
//            // 左移窗口
//            left++;
//            // 进行窗口内数据的一系列更新
//            ...
//        }
//    }
// }

// 给你一个字符串 S、一个字符串 T，请在字符串 S 里面找出：包含 T 所有字母的最小子串
func minWindow(s string, t string) string {
	// 保存滑动窗口字符集
	win := make(map[byte]int)
	// 保存需要的字符集
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	// 窗口
	left := 0
	right := 0
	// match匹配次数
	match := 0
	start := 0
	end := 0
	min := math.MaxInt64
	var c byte
	for right < len(s) {
		c = s[right]
		right++
		// 在需要的字符集里面，添加到窗口字符集里面
		if need[c] != 0 {
			win[c]++
			// 如果当前字符的数量匹配需要的字符的数量，则match值+1
			if win[c] == need[c] {
				match++
			}
		}

		// 当所有字符数量都匹配之后，开始缩紧窗口
		for match == len(need) {
			// 获取结果
			if right-left < min {
				min = right - left
				start = left
				end = right
			}
			c = s[left]
			left++
			// 左指针指向不在需要字符集则直接跳过
			if need[c] != 0 {
				// 左指针指向字符数量和需要的字符相等时，右移之后match值就不匹配则减一
				// 因为win里面的字符数可能比较多，如有10个A，但需要的字符数量可能为3
				// 所以在压死骆驼的最后一根稻草时，match才减一，这时候才跳出循环
				if win[c] == need[c] {
					match--
				}
				win[c]--
			}
		}
	}
	if min == math.MaxInt64 {
		return ""
	}
	return s[start:end]
}
