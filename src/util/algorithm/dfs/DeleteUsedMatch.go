package dfs

import (
	"fmt"
	"sort"
)

// 删除最小数量的无效括号，使得输入的字符串有效，返回所有可能的结果。
//
//说明: 输入可能包含了除 ( 和 ) 以外的字符。
//
//示例 1:
//
//输入: "()())()"
//输出: ["()()()", "(())()"]
//示例 2:
//
//输入: "(a)())()"
//输出: ["(a)()()", "(a())()"]
//示例 3:
//
//输入: ")("
//输出: [""]
//

// 从第一个字符开始每一个字符除了非括号类型都有在/或者不在结果集合里面的可能性，对于非括号类型，默认是在结果集合里面的。
// 我们需要做的就是从第一个字符开始，采用dfs深度优先遍历的访问方式，遍历到最后一个字符看是否合法。
// 具体实现上，我们可以一遍遍历，一边记录count的情况，--即左右括号的匹配情况。。。
// 在最后的一个字符遍历完成后，只需要判断左右括号是否为0即可。
// 怎么保证删掉最少的括号？
//
//这个方法很多，说一下我的。假设我们用 res 数组保存最终的结果，当新的字符串要加入的时候，我们判断一下新加入的字符串的长度和数组中第一个元素长度的关系。
//
//如果新加入的字符串的长度大于数组中第一个元素的长度，我们就清空数组，然后再将新字符串加入。
//
//如果新加入的字符串的长度小于数组中第一个元素的长度，那么当前字符串抛弃掉。
//
//如果新加入的字符串的长度等于数组中第一个元素的长度，将新字符串加入到 res 中。
//
//第四，重复的情况怎么办？
//
//简单粗暴一些，最后通过 set 去重即可。

func removeInvalidParentheses(s string) []string {
	temp_res := make([]string, 0)
	temp_res = append(temp_res, "")
	if len(s) == 0 {
		return temp_res
	}
	dfs(s, 0, len(s), "", &temp_res, 0)
	// 去重
	dict := make(map[string]bool)
	res := make([]string, 0)
	for index := range temp_res {
		if _, ok := dict[temp_res[index]]; ok {
			continue
		} else {
			dict[temp_res[index]] = true
			res = append(res, temp_res[index])
		}
	}
	return res
}

func dfs(s string, start int, length int, temp string, res *[]string, count int) {
	if count < 0 {
		return // 右边括号太多，成为合法的括号的可能性很小了，byte byte
	}
	// 到达字符串末尾
	if start == length {
		if count == 0 {
			// valid
			currMax := 0

			if len(*res) > 0 {
				currMax = len((*res)[0])
			}
			if len(temp) > currMax {
				// 清空历史
				*res = (*res)[0:0]        // clear res
				*res = append(*res, temp) // 加入当前元素

			} else if currMax == len(temp) {
				*res = append(*res, temp)
			}
		}
		return
	}
	// 即任何一个字符都存在加入或者不加入两种选择！！！
	// 添加当前字符的场景 add or not 主体的思想考虑  if or not !!! todo
	if s[start] == '(' {
		dfs(s, start+1, length, temp+"(", res, count+1)
	} else if s[start] == ')' {
		dfs(s, start+1, length, temp+")", res, count-1)
	} else {
		// 非括号字符都是默认添加的
		dfs(s, start+1, length, temp+s[start:start+1], res, count)
	}
	// 不添加当前字符 -- 主要是针对括号类型字符是否进行加入考虑的
	if s[start] == '(' || s[start] == ')' {
		// 直接偏移量 + 1， 临时结果集合里面不需要加上该字符， count保持不变
		dfs(s, start+1, length, temp, res, count)
	}
}

// // 思考：记dp[i] 表示 s[0...i]的最长有效字符串的长度, / 以str[i]结尾的字串最大可以表示的合法括号字符串长度
//// if s[i] == ')' ,则分析 以 s[i-1]结尾的字符串所能够形成的最大的最大的长度，排除掉这一部分字符的长度后，
//// 还要考虑前面可能构成的最大字符串长度
func longestValidQuote(str string) ([]int, int) {
	dp := make([]int, len(str))
	if len(str) == 0 {
		return dp, 0
	}
	res := 0
	for i := 1; i < len(str); i++ {
		if str[i] != '(' && str[i] != ')' {
			dp[i] = dp[i] + 1
		}
		if str[i] == ')' {
			// i - 1 - x  + 1 = dp[i-1] => dp[i-1] = i - x
			offset := i - dp[i-1] - 1 // 前一个字符匹配的最长有效括号的前一个字符
			if offset >= 0 {
				if str[offset] == '(' {
					// 考虑假如匹配的情况下，匹配的前面部分可能也会得以匹配
					pre := offset - 1
					dp[i] = dp[i-1] + 2
					if pre >= 0 {
						dp[i] = dp[i] + dp[pre]
					}
				}
			}
		}
		if res < dp[i] {
			res = dp[i]
		}
	}
	return dp, res
}

// 判断括号序列是否是合法的匹配序列
func isValid(s string) bool {
	if len(s) == 0 {
		return false
	}
	count := 0
	for _, e := range s {
		if e == '(' {
			count += 1
		} else if e == ')' {
			count -= 1
		}
		if count < 0 { // 右边括号显得过于多了
			return false
		}
	}
	return count == 0
}

func removeInvalidParentheses2(s string) []string {
	res := make([]string, 0)
	//res = append(res, "")
	if len(s) == 0 {
		return res
	}
	bfs(&res, s)
	if len(res) == 0 {
		res = append(res, "")
		return res
	}
	sort.Slice(res, func(i, j int) bool {
		return len(res[i]) > len(res[j])
	})
	return res
}

// bfs的代码看起来显得更加简洁一点儿
func bfs(res *[]string, s string) {
	levels := make(map[string]int)
	levels[s] = 0 // 最顶层加入原始的字符集合 ！！！
	for {
		for k := range levels {
			if isValid(k) {
				*res = append(*res, k)
			}
		}
		if len(*res) > 0 {
			return
		}
		nextLevels := make(map[string]int) // 生成下一层新的level ，进行bfs操作
		for item := range levels {
			if item == "(" || item == ")" {
				return
			}
			for i := 0; i < len(item); i++ {
				if item[i] == ')' || item[i] == '(' {
					e := item[:i] + item[i+1:] // 去除掉当前字符！！
					if _, ok := nextLevels[e]; !ok {
						nextLevels[e] = 0
					}
				}
			}
		}
		levels = nextLevels
	}
}

func TestLongestLenght(s string) {
	res := removeInvalidParentheses2(s)
	fmt.Println("res is :", res)
}
