package tree

// 线段树(Segment Tree)也是一棵树，只不过元素的值代表一个区间。常用区间的统计操作，
// 比如一个区间的最大值(max)，最小值(min)，和(sum)等等，线段树是一个平衡二叉树，但不一定是完全二叉树。

// 根节点就是 0~lenght-1 的和，根节点的左右子树平分根节点的区间，然后依次类推，直到只有一个元素不能划分为止，该元素也就是二叉树的叶子节点。
//
//求线段树的区间统计，时间复杂度和二叉树的高度有关系，和元素的个数没关系，它的时间复杂度为 O(log n)。

// h层的满二叉树总共有 2^h - 1个节点，地h-1层具有 2^(h-1)个节点，他们大概是2倍的关系，也就是说满二叉树最后一层的节点数目乘以
// 2大概就是总的节点数目，但是线段树不是满二叉树，但是一定是平衡二叉树，所以需要冗余一层

// n个元素构成的线段树的叶子节点总共具有n个，n个叶子节点的完全二叉树具有n-1个非叶子节点
// 所以总共的节点数目为：2n -1,  这2n-1个只是有效的节点，因为我们用数组存储，所以需要给这颗近似的完全二叉树添加一些虚节点，让其变成一颗真正的完全二叉树。
//对于n个元素的序列构成的线段树，2x-1（x是大于等于n的最小的2的幂，比如n=6则x=8）个节点绝对够用，因为n个元素的序列构建的完全二叉树补全成为满二叉树，该满二叉树的节点个数必为2x-1。

// 求线段树的区间统计[和 、最大、最小]，时间复杂度和二叉树的高度有关系，和元素的个数没关系，它的时间复杂度为 O(log n)。

type Node struct {
	left  int // 左边界
	right int // 右边界
	data  int // 节点数值：权重
	mark  int // mark标记，方便延迟更新操作
}

// 线段树
type SegmentTree struct {
	nums  []int
	nodes []Node // 节点表示
}

// 构建线段树结构
func NewSegmentTree(nums []int) SegmentTree {
	segmentTree := SegmentTree{nums: nums}
	if nums == nil || len(nums) == 0 {
		return segmentTree
	}
	segmentTree.nodes = make([]Node, 2*len(nums)+1) // 初始化nodes节点
	segmentTree.buildSegmentTree(0)
	//fmt.Println("xxxx")
	return segmentTree
}

func (this *SegmentTree) buildSegmentTree(index int) {
	node := this.nodes[index]
	if node.mark == 0 {
		// not not exist
		node = Node{left: 0, right: len(this.nums) - 1}
		this.nodes[index] = node
		//node = this.nodes[index]
		node.mark = 1
	}
	if node.left == node.right {
		// 🍃节点
		//fmt.Println("node data is ", this.nums[node.left])
		node.data = this.nums[node.left]
		this.nodes[index] = node
	} else {
		mid := (node.left + node.right) >> 1
		this.nodes[(index<<1)+1] = Node{left: node.left, right: mid, mark: 1}      // 左节点线段
		this.nodes[(index<<1)+2] = Node{left: mid + 1, right: node.right, mark: 1} // 右节点线段
		this.buildSegmentTree((index << 1) + 1)                                    // 递归构建左子树
		this.buildSegmentTree((index << 1) + 2)                                    // 递归构建右子树
		this.nodes[index].data = this.nodes[(index<<1)+1].data + this.nodes[index<<1+2].data
	}
	// fmt.Println("nodes is :", this.nodes)
	//node.mark = 1 // 表示node被初始化完成
}

// 查询指定区间的和
func (this *SegmentTree) query(index, left, right int) int {
	node := this.nodes[index]
	if node.left == node.left && node.right == right {
		// 当前区间和带查询区间完全匹配
		return this.nodes[index].data
	}
	mid := (node.left + node.right) >> 1 // 区间中点,对应左孩子区间结束,右孩子区间开头
	if right <= mid {
		// 查询区间全在左子树
		return this.query(index<<1+1, left, right)
	}
	if left > mid {
		return this.query(index<<1+2, left, right)
	}
	return this.query(index<<1+1, left, right) + this.query(index<<1+2, left, right)
}

// index 根节点index
// updateIndex 待更新节点的下标
// left 、right 跟节点左右边界
// 思路：先定位到叶子节点更新值，然后更新包含该节点的区间的值, 直到根节点为止
func (this *SegmentTree) updateSegment(index, updateIndex, left, right int) {
	node := this.nodes[index]
	if node.left == node.right && node.left == updateIndex {
		this.nodes[index].data = this.nums[updateIndex] // update val
		return
	}
	mid := (node.left + node.right) >> 1
	l, r := index<<1+1, index<<1+2
	if updateIndex > mid {
		// 更新的区间在右边子树
		this.updateSegment(r, updateIndex, mid+1, node.right)
	}
	if updateIndex <= mid {
		this.updateSegment(l, updateIndex, node.left, mid)
	}
	// 更新父节点值
	this.nodes[index].data = this.nodes[l].data + this.nodes[r].data
}

type NumArray struct {
	nums        []int
	segmentTree SegmentTree
}

func Constructor(nums []int) NumArray {

	return NumArray{
		nums:        nums,
		segmentTree: NewSegmentTree(nums),
	}
}

func (this *NumArray) Update(i int, val int) {
	if len(this.nums) == 0 {
		return
	}
	if i >= len(this.nums) {
		return
	}
	this.nums[i] = val
	this.segmentTree.updateSegment(0, i, 0, len(this.nums)-1)
}

func (this *NumArray) SumRange(i int, j int) int {
	if len(this.nums) == 0 {
		return 0
	}

	return this.segmentTree.query(0, i, j)
}
