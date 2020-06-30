package consientHash

import (
	"crypto/sha1"
	"math"
	"sort"
	"strconv"
	"sync"
)

const (
	//DefaultVirualSpots default virual spots
	DefaultVirualSpots = 400
)

// 服务器机器节点
type node struct {
	nodeKey   string
	spotValue uint32
}

type nodesArray []node

// 实现了Sort接口，可以用于比较
func (p nodesArray) Len() int           { return len(p) }
func (p nodesArray) Less(i, j int) bool { return p[i].spotValue < p[j].spotValue }
func (p nodesArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p nodesArray) Sort()              { sort.Sort(p) }

//HashRing store nodes and weigths 2^32
// 这里存在一种场景, 当一个服务由多个服务器组共同提供时, key应该路由到哪一个服务.这里假如采用最通用的方式key%N(N为服务器数目),
//这里乍一看没什么问题, 但是当服务器数目发送增加或减少时,
//分配方式则变为key%(N+1)或key%(N-1).这里将会有大量的key失效迁移,如果后端key对应的是有状态的存储数据
//,那么毫无疑问,这种做法将导致服务器间大量的数据迁移,从而照成服务的不稳定. 为了解决类问题,一致性hash算法应运而生.

// 一致性hash算法特点
//在分布式缓存中, 一个好的hash算法应该要满足以下几个条件:
//均衡性(Balance)
//均衡性主要指,通过算法分配, 集群中各节点应该要尽可能均衡.
//单调性(Monotonicity)
//单调性主要指当集群发生变化时, 已经分配到老节点的key, 尽可能的任然分配到之前节点,以防止大量数据迁移, 这里一般的hash取模就很难满足这点,
//而一致性hash算法能够将发生迁移的key数量控制在较低的水平.
//分散性(Spread)
//分散性主要针对同一个key, 当在不同客户端操作时,可能存在客户端获取到的缓存集群的数量不一致,
//从而导致将key映射到不同节点的问题,这会引发数据的不一致性.好的hash算法应该要尽可能避免分散性.
//
//负载(Load)
//负载主要是针对一个缓存而言, 同一缓存有可能会被用户映射到不同的key上,从而导致该缓存的状态不一致.
//
//从原理来看,一致性hash算法针对以上问题均有一个合理的解决.

// 一致性hash的核心思想为将key作hash运算, 并按一定规律取整得出0-2^32-1之间的值, 环的大小为2^32，key计算出来的整数值则为key在hash环上的位置，如何将一个key，映射到一个节点， 这里分为两步.
//第一步, 将服务的key按该hash算法计算,得到在服务在一致性hash环上的位置.
//第二步, 将缓存的key，用同样的方法计算出hash环上的位置，按顺时针方向，找到第一个大于等于该hash环位置的服务key，从而得到该key需要分配的服务器。
// 如图， 各key根据hash算法分配到各节点，当某一节点失效实效时， 如NODE 2， 则NODE 2 上的key将分配到hash环上相邻的节点，而其他key所在位置不变。
// 虚拟节点提高均衡性
//
//如上图可看到， 由于节点只有3个，存在某些节点所在位置周围有大量的hash点从而导致分配到这些节点到key要比其他节点多的多，这样会导致集群中各节点负载不均衡，
//为解决这个问题，引入虚拟节点， 即一个实节点对应多个虚拟节点。缓存的key作映射时，先找到对应的虚拟节点，再对应到实节点。如下图所示， 每个节点虚拟出两个虚拟节点，从而提高均衡性。
type HashRing struct {
	virualSpots int            // 虚拟节点的个数
	nodes       nodesArray     // 具体的节点
	weights     map[string]int // 每个节点都有一个权重
	mu          sync.RWMutex   // 读写锁
}

//NewHashRing create a hash ring with virual spots
func NewHashRing(spots int) *HashRing {
	if spots == 0 {
		spots = DefaultVirualSpots
	}

	h := &HashRing{
		virualSpots: spots,
		weights:     make(map[string]int),
	}
	return h
}

//AddNodes add nodes to hash ring
func (h *HashRing) AddNodes(nodeWeight map[string]int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for nodeKey, w := range nodeWeight {
		h.weights[nodeKey] = w
	}
	h.generate()
}

//AddNode add node to hash ring
func (h *HashRing) AddNode(nodeKey string, weight int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[nodeKey] = weight
	h.generate()
}

//RemoveNode remove node
func (h *HashRing) RemoveNode(nodeKey string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.weights, nodeKey)
	h.generate()
}

//UpdateNode update node with weight
func (h *HashRing) UpdateNode(nodeKey string, weight int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[nodeKey] = weight
	h.generate()
}

// 重新修正映射关系。。。
func (h *HashRing) generate() {
	var totalW int // 总权重值
	for _, w := range h.weights {
		totalW += w
	}
	// 总共的虚拟节点的数目
	totalVirtualSpots := h.virualSpots * len(h.weights)
	h.nodes = nodesArray{}

	for nodeKey, w := range h.weights {
		// 向上取整，权重越大的节点，的虚拟节点数目也对应的越大。。
		spots := int(math.Floor(float64(w) / float64(totalW) * float64(totalVirtualSpots)))
		for i := 1; i <= spots; i++ {
			hash := sha1.New()
			hash.Write([]byte(nodeKey + ":" + strconv.Itoa(i)))
			hashBytes := hash.Sum(nil)
			n := node{
				nodeKey:   nodeKey,
				spotValue: genValue(hashBytes[6:10]),
			}
			h.nodes = append(h.nodes, n)
			hash.Reset()
		}
	}
	h.nodes.Sort() // 排序，方便请求key顺时针搜索定位
}

// 在2^32的环上映射slot!!!
func genValue(bs []byte) uint32 {
	if len(bs) < 4 {
		return 0
	}
	v := (uint32(bs[3]) << 24) | (uint32(bs[2]) << 16) | (uint32(bs[1]) << 8) | (uint32(bs[0]))
	return v
}

//GetNode get node with key
func (h *HashRing) GetNode(s string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.nodes) == 0 {
		return ""
	}

	hash := sha1.New()
	hash.Write([]byte(s))
	hashBytes := hash.Sum(nil)
	v := genValue(hashBytes[6:10])
	i := sort.Search(len(h.nodes), func(i int) bool { return h.nodes[i].spotValue >= v })

	if i == len(h.nodes) {
		i = 0
	}
	return h.nodes[i].nodeKey
}
