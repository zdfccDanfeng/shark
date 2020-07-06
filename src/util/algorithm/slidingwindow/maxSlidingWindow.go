package slidingwindow

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
