package nameserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc/resolver"
	"strings"
	"time"
)

// 整个自定义命名解析功能最核心的代码，通过自定义 MyCustomResolver 将服务名解析成对应实例
const (
	// syncNSInterval 定义了从 NS 服务同步实例列表的周期
	syncNSInterval = 1 * time.Second
)

type MyCustomResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

// ResolveNow 实现了 resolver.Resolver.ResolveNow 方法
func (this *MyCustomResolver) ResolveNow(o resolver.ResolveNowOptions) {
	this.watcher()
}

// watcher 轮询并更新指定 CalleeService 服务的实例变化
func (r *MyCustomResolver) watcher() {
	r.updateCC()
	ticker := time.NewTicker(syncNSInterval)
	for {
		select {
		// 当* nsResolver Close 时退出监听
		case <-r.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			// 调用* nsResolver.updagteCC() 方法，更新实例地址
			r.updateCC()
		}
	}
}

// instances 包含调用方服务名、被调方服务名、被调方实例列表等数据
type instances struct {
	callerService string
	calleeService string
	calleeIns     []string
	CalleeIns     resolver.Address
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
		calleeIns:     ins.Instances,
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
		newAddrs = append(newAddrs, instances.CalleeIns)
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

type DnsInfo struct {
	Instances []string // 解析到的被调地址列表信息
}

func GetInstances(auth string, endPoint string) (DnsInfo, error) {
	return DnsInfo{}, nil
}
