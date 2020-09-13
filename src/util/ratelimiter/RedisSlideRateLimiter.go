package ratelimiter

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

//前面我们提到限流的主要目的是为了保证系统的稳定性。在日常的业务中，如果遇到像双十一之类的促销活动，或者遇到爬虫等不正常的流量等情况，用户流量突增，但后端服务的处理能力是有限的，如果不能有效处理突发流量，那么后端服务就很容易被打垮。
//可以设想这样一个场景："某服务单节点可以承受的QPS是1000，该服务共有5个节点，日常情况下服务的QPS为3000"。那么正常情况下该服务毫无压力,根据负载均衡配置3000/5=600，每个节点的日常QPS才600左右。
//直到某一天，老板突然搞了一波促销，系统的整体QPS达到了8000。此时每个节点的平均承载QPS为1600，节点A率先扛不住直接挂了，此时集群中还剩下4个节点，每个节点的平均承载QPS将达到2000，于是，剩下的4个节点也一台接一台挂了，整个服务就此雪崩。
//而如果我们的系统有限流机制，那么情况将会如何发展呢？
//系统整体QPS达到8000，但由于集群整体限流了5000，所以超出集群承受力的那3000个请求将被拒绝，系统则会正常处理5000个用户请求，这是对集群整体限流的情况。而对于各个节点来说，由于我的承受力只是1000QPS，那么超出1000的部分也将被拒绝。
//这样虽然损失了部分用户请求，但保证了整个系统的稳定性，也给开发运维留下了系统扩容时间。
// 常见的限流算法主要有：计数器、固定窗口，滑动窗口、漏桶、令牌桶。接下来我们分别介绍下这几种限流算法。
func DoRateLimiter() {
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1",
		Password: "",
		DB:       0})

	for i := 0; i < 20; i++ {
		//假设对用户jankin的登录操作进行限流检测，60秒内允许登录5次
		fmt.Println(isActionAllowed(*rdb, "jankin", "login", 60, 5))
		//可以根据isActionAllowed方法返回的是true还是false来判断是否达到限流阈值
	}
}

// 当然是用Redis限流的主流方法还有漏桶算法（leaky-bucket）和令牌桶算法（token-bucket），本文主要讲解简单的计数器和滑动窗口法，这两种算法都属于计数器法，后面将有更详细的实验进行介绍漏桶算法和令牌桶算法。
func isActionAllowed(rdb redis.Client, userId, actionKey string, period, maxCount int) bool {
	millisecond := time.Now().UnixNano() / 1e6
	key := userId + "_" + actionKey
	//裁剪zset，只保留窗口内的值
	rdb.ZRemRangeByScore(key, "0", strconv.FormatInt(millisecond-int64(period*1000), 10))
	//统计窗口的值的数量
	ran := rdb.ZCard(key)
	//先判断是否达到阈值
	if int(ran.Val()) < maxCount {
		fmt.Println("当前未到达限流阈值,插入新数据：", millisecond)
		//符合阈值才添加数据
		rdb.ZAdd(key, redis.Z{
			Score:  float64(millisecond),
			Member: millisecond,
		})
		//设置一个过期时间，通常只要比窗口大一点即可，这里设置大1秒
		rdb.Expire(key, time.Duration(period+1)*time.Second)
		return true
	} else {
		fmt.Println("当前已到达限流阈值")
		return false
	}
}
