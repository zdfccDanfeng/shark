package tree

import (
	"fmt"
	"strconv"
	"strings"
)

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// 给定一个二叉树，返回所有从根节点到叶子节点的路径。
//
//说明: 叶子节点是指没有子节点的节点。
//
//示例:
//
//输入:
//
//   1
// /   \
//2     3
// \
//  5
//
//输出: ["1->2->5", "1->3"]
//
//解释: 所有根节点到叶子节点的路径为: 1->2->5, 1->3
//
func binaryTreePaths(root *TreeNode) []string {
	res := make([]string, 0)
	if root == nil {
		return res
	}
	temp := make([]string, 0)
	dfs(root, &res, temp)
	return res
}

func dfs(root *TreeNode, res *[]string, temp []string) {
	if root == nil {
		return
	}

	// 前序列遍历
	temp = append(temp, strconv.Itoa(root.Val))
	dfs(root.Left, res, temp)
	dfs(root.Right, res, temp)
	if root.Right == nil && root.Left == nil {
		*res = append(*res, strings.Join(temp, "->"))
	}
	temp = temp[:len(temp)-1] // 去掉最后一个
}

func TestBinaryTreePaths(root *TreeNode) {
	paths := binaryTreePaths(root)
	fmt.Println("res is :", paths)
}
