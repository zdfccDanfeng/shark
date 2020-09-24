package combination

// 找出所有相加之和为 n 的 k 个数的组合。组合中只允许含有 1 - 9 的正整数，并且每种组合中不存在重复的数字。
// 说明：
// 所有数字都是正整数。
// 解集不能包含重复的组合。
// 示例 1:
//
//输入: k = 3, n = 7
//输出: [[1,2,4]]
//示例 2:
//
//输入: k = 3, n = 9
//输出: [[1,2,6], [1,3,5], [2,3,4]]
func combinationSum3(k int, n int) [][]int {
	res := make([][]int, 0)
	temp := make([]int, 0)
	dfs(n, k, 0, 1, 0, &res, temp)
	return res
}

func dfs(n, k, c, start, current int, res *[][]int, temp []int) {
	if c == k {
		if current == n {
			*res = append(*res, temp[:])
		}
		return
	}
	for t := start; t <= 9; t++ {
		temp = append(temp, t)
		dfs(n, k, c+1, t+1, current+t, res, temp)
		temp = temp[:len(temp)-1]
	}
}
