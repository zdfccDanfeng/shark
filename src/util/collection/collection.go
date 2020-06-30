package collection

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/treeset"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"github.com/emirpasic/gods/utils"
)

func testList() {
	list := sll.New()
	list.Add("a")                         // ["a"]
	list.Add("c", "b")                    // ["a","c","b"]
	list.Sort(utils.StringComparator)     // ["a","b","c"]
	_, _ = list.Get(0)                    // "a",true
	_, _ = list.Get(100)                  // nil,false
	_ = list.Contains("a", "b", "c")      // true
	_ = list.Contains("a", "b", "c", "d") // false
	list.Swap(0, 1)                       // ["b","a",c"]
	list.Remove(2)                        // ["b","a"]
	list.Remove(1)                        // ["b"]
	list.Remove(0)                        // []
	list.Remove(0)                        // [] (ignored)
	_ = list.Empty()                      // true
	_ = list.Size()                       // 0
	list.Add("a")                         // ["a"]
	list.Clear()                          // []
	list.Insert(0, "b")                   // ["b"]
	list.Insert(0, "a")                   // ["a","b"]

	// ==================================
	list2 := dll.New()
	list2.Add("a")                        // ["a"]
	list2.Add("c", "b")                   // ["a","c","b"]
	list.Sort(utils.StringComparator)     // ["a","b","c"]
	_, _ = list.Get(0)                    // "a",true
	_, _ = list.Get(100)                  // nil,false
	_ = list.Contains("a", "b", "c")      // true
	_ = list.Contains("a", "b", "c", "d") // false
	list.Swap(0, 1)                       // ["b","a",c"]
	list.Remove(2)                        // ["b","a"]
	list.Remove(1)                        // ["b"]
	list.Remove(0)                        // []
	list.Remove(0)                        // [] (ignored)
	_ = list.Empty()                      // true
	_ = list.Size()                       // 0
	list.Add("a")                         // ["a"]
	list.Clear()                          // []
	list.Insert(0, "b")                   // ["b"]
	list.Insert(0, "a")                   // ["a","b"]

	// =============
	set := hashset.New()   // empty
	set.Add(1)             // 1
	set.Add(2, 2, 3, 4, 5) // 3, 1, 2, 4, 5 (random order, duplicates ignored)
	set.Remove(4)          // 5, 3, 2, 1 (random order)
	set.Remove(2, 3)       // 1, 5 (random order)
	set.Contains(1)        // true
	set.Contains(1, 5)     // true
	set.Contains(1, 6)     // false
	_ = set.Values()       // []int{5,1} (random order)
	set.Clear()            // empty
	set.Empty()            // true
	set.Size()             // 0

	// ========

	set3 := treeset.NewWithIntComparator() // empty (keys are of type int)
	set3.Add(1)                            // 1
	set3.Add(2, 2, 3, 4, 5)                // 1, 2, 3, 4, 5 (in order, duplicates ignored)
	set.Remove(4)                          // 1, 2, 3, 5 (in order)
	set.Remove(2, 3)                       // 1, 5 (in order)
	set.Contains(1)                        // true
	set.Contains(1, 5)                     // true
	set.Contains(1, 6)                     // false
	_ = set.Values()                       // []int{1,5} (in order)
	set.Clear()                            // empty
	set.Empty()                            // true
	set.Size()                             // 0

	// ==================
	stack := lls.New()  // empty
	stack.Push(1)       // 1
	stack.Push(2)       // 1, 2
	stack.Values()      // 2, 1 (LIFO order)
	_, _ = stack.Peek() // 2,true
	_, _ = stack.Pop()  // 2, true
	_, _ = stack.Pop()  // 1, true
	_, _ = stack.Pop()  // nil, false (nothing to pop)
	stack.Push(1)       // 1
	stack.Clear()       // empty
	stack.Empty()       // true
	stack.Size()        // 0

	m := hashmap.New() // empty
	m.Put(1, "x")      // 1->x
	m.Put(2, "b")      // 2->b, 1->x (random order)
	m.Put(1, "a")      // 2->b, 1->a (random order)
	_, _ = m.Get(2)    // b, true
	_, _ = m.Get(3)    // nil, false
	_ = m.Values()     // []interface {}{"b", "a"} (random order)
	_ = m.Keys()       // []interface {}{1, 2} (random order)
	m.Remove(1)        // 2->b
	m.Clear()          // empty
	m.Empty()          // true
	m.Size()           // 0

	m2 := treemap.NewWithIntComparator() // empty (keys are of type int)
	m2.Put(1, "x")                       // 1->x
	m.Put(2, "b")                        // 1->x, 2->b (in order)
	m.Put(1, "a")                        // 1->a, 2->b (in order)
	_, _ = m.Get(2)                      // b, true
	_, _ = m.Get(3)                      // nil, false
	_ = m.Values()                       // []interface {}{"a", "b"} (in order)
	_ = m.Keys()                         // []interface {}{1, 2} (in order)
	m.Remove(1)                          // 2->b
	m.Clear()                            // empty
	m.Empty()                            // true
	m.Size()                             // 0

	// Other:
	m2.Min() // Returns the minimum key and its value from map.
	m2.Max() // Returns the maximum key and its value from map.

	m3 := linkedhashmap.New() // empty (keys are of type int)
	m3.Put(2, "b")            // 2->b
	m.Put(1, "x")             // 2->b, 1->x (insertion-order)
	m.Put(1, "a")             // 2->b, 1->a (insertion-order)
	_, _ = m.Get(2)           // b, true
	_, _ = m.Get(3)           // nil, false
	_ = m.Values()            // []interface {}{"b", "a"} (insertion-order)
	_ = m.Keys()              // []interface {}{2, 1} (insertion-order)
	m.Remove(1)               // 2->b
	m.Clear()                 // empty
	m.Empty()                 // true
	m.Size()                  // 0

}
