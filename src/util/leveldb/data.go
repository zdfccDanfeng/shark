package leveldb

import (
	"github.com/emirpasic/gods/utils"
	"sync"
)

type ValueType int8

// 使用Type字段来表明这是删除or插入一个新Key
const (
	TypeDeletion ValueType = 1
	TypeValue    ValueType = 2
)

type InternalKey struct {
	Seq       int64     // seq number,原子自增
	Type      ValueType // add or delete
	UserKey   []byte    // real key
	UserValue []byte    // real value [delete with empty]
}

const (
	kMaxHeight = 12
	kBranching = 4
)

type Node struct {
	key  interface{}
	next []*Node
}
type SkipList struct {
	maxHeight  int
	head       *Node
	comparator utils.Comparator
	locker     sync.RWMutex
}

type MemTable struct {
	table       *SkipList
	memoryUsage uint64
}
