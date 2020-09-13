package nameserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc/resolver"
	"strings"
)

type MyCustomResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func (this *MyCustomResolver) ResolveNow(o resolver.ResolveNowOptions) {

}

// instances 包含调用方服务名、被调方服务名、被调方实例列表等数据
type instances struct {
	callerService string
	calleeService string
	calleeIns     []string
}

func getInstances(target resolver.Target) (s *instances, e error) {
	auths := strings.Split(target.Authority, "@")
	// auths[0] 为 callerService 名，target.Endpoint 为 calleeService 名
	// 通过自定义 sdk 从内部 NameServer 查询指定 calleeService 对应的实例列表
	ins, e := GetInstances(auths[0], target.Endpoint)
	if e != nil {
		return nil, e
	}
	return &instances{
		callerService: auths[0],
		calleeService: target.Endpoint,
		calleeIns:     ins,
	}, nil

}

// updateCC 更新 resolver.Resolver.ClientConn 配置
func (r *MyCustomResolver) updateCC() {
	// 从 NS 服务获取指定 target 的实例列表
	instances, err := getInstances(r.target)
	// 如果获取实例列表失败，或者实例列表为空，则不更新 resolver 中实例列表
	if err != nil || len(instances.calleeIns) == 0 {
		//logz.Warn("[mis] error retrieving instances from Mis", logz.Any("target", r.target), logz.Error(err))
		return
	}
	//...

	// 组装实例列表 []resolver.Address
	// resolver.Address 结构体表示 grpc server 端实例地址
	var newAddrs []resolver.Address
	for k := range instances.calleeIns {
		//	newAddrs = append(newAddrs, instances.calleeIns)
		fmt.Println(k)
	}
	//...

	// 更新实例列表
	// grpc 底层 LB 组件对每个服务端实例创建一个 subConnection。并根据设定的 LB 策略，选择合适的 subConnection 处理某次 RPC 请求。
	// 此处代码比较复杂，后续在 LB 相关原理文章中再做概述
	r.cc.UpdateState(resolver.State{Addresses: newAddrs})
}

func (this *MyCustomResolver) Close() {
	this.cancel()
}

func GetInstances(auth string, endPoint string) ([]string, error) {
	return []string{}, nil
}
