package dp

import "math"

// 给定一个非负整数数组和一个整数 m，你需要将这个数组分成 m 个非空的连续子数组。设计一个算法使得这 m 个子数组各自和的最大值最小。
//
//注意:
//数组长度 n 满足以下条件:
//
//1 ≤ n ≤ 1000
//1 ≤ m ≤ min(50, n)
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/split-array-largest-sum
//著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
func SplitArray(nums []int, m int) int {

	n := len(nums)
	var dp [1001][51]int
	for i := 0; i < len(dp); i++ {
		for j := 0; j < len(dp[0]); j++ {
			dp[i][j] = math.MaxInt32
		}
	}
	dp[0][0] = math.MaxInt32

	// 记dp[i][j] 表示 数组的前i个序列，被划分成j组情况下的最小最大值，
	// 那么 i从1到k可以被划分成j-1个序列，k+1到i可以被划分成第j个序列。todo
	//则可以得到如下的递推公式：
	// dp[i][j] = min{dp[i][j] , max(dp[k][j-1], (sum[i] - sum[k])}, 其中 1<=k< i
	var sum [1001]int
	sum[0] = 0
	for i := 1; i <= n; i++ {
		sum[i] = sum[i-1] + nums[i-1]
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= min(i, m); j++ {
			if j == 1 {
				dp[i][j] = sum[i]
			} else {
				for k := 1; k < i; k++ {
					dp[i][j] = min(dp[i][j], max(dp[k][j-1], sum[i]-sum[k]))
				}
			}
		}
	}
	return dp[n][m]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
