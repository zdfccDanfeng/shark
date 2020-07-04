package list

import "fmt"

type ListNode struct {
	Value int
	Next  *ListNode
}

func NewListNode(val int) *ListNode {
	return &ListNode{Value: val, Next: nil}
}

// 1---->2----->3 递归反转单链表操作
// 1  3--->2
//
func SwapListNode(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	currNode := head
	preNode := SwapListNode(currNode.Next) // 对后半部分先进行反转操作
	head.Next.Next = currNode              // 巧妙利用递归返回后，头指针始终指向剩余反转后链表的尾部节点。
	currNode.Next = nil
	return preNode
}

func (list *ListNode) PrintNode() {
	curr := list
	for {
		if curr == nil {
			break
		}
		fmt.Println(curr.Value)
		curr = curr.Next
	}
}
