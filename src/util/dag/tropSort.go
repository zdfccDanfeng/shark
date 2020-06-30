package dag

import (
	"fmt"
	"sort"
)

/**
拓扑排序
从概念上说，前置条件可以构成有向图。图中的顶点表示课程，边表示课程间的依赖关系。
显然，图中应该无环，这也就是说从某点出发的边，
最终不会回到该点。下面的代码用深度优先搜索了整张图，获得了符合要求的课程序列
dfs思想 --》 拓扑排序
*/
func DagSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}

func TestDAg(m map[string][]string) {
	res := DagSort(m)
	for _, e := range res {
		fmt.Println(e)
	}
}
