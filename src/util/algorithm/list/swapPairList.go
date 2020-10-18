package list

//* Definition for singly-linked list.

//给定一个链表，两两交换其中相邻的节点，并返回交换后的链表。
//你不能只是单纯的改变节点内部的值，而是需要实际的进行节点交换。
// 思想：滚动更新方式。。。
func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	// 增加虚拟头节点，避免判空
	dumpy := &ListNode{Value: -1, Next: head}
	pre := dumpy
	for head != nil && head.Next != nil {
		first := head
		second := head.Next
		pre.Next = second
		first.Next = second.Next
		second.Next = first
		pre = first
		head = pre.Next
	}

	return dumpy.Next
}
