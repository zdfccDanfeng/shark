package zookeeper

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
)

// grcp提供了进程内负载均衡的设计思路
// gRPC开源组件官方并未直接提供服务注册与发现的功能实现，但其设计文档已提供实现的思路，并在不同语言的gRPC代码API中已提供了命名解析和负载均衡接口供扩展
// https://blog.csdn.net/Edu_enth/article/details/104016307
// 其基本实现原理：
//1。 服务启动后gRPC客户端向命名服务器发出名称解析请求，名称将解析为一个或多个IP地址，每个IP地址标示它是服务器地址还是负载均衡器地址，
//  以及标示要使用那个客户端负载均衡策略或服务配置。
//2。 客户端实例化负载均衡策略，如果解析返回的地址是负载均衡器地址，则客户端将使用grpclb策略，否则客户端使用服务配置请求的负载均衡策略。
//3。 负载均衡策略为每个服务器地址创建一个子通道（channel）。
//4。 当有rpc请求时，负载均衡策略决定那个子通道即grpc服务器将接收请求，当可用服务器为空时客户端的请求将被阻塞。

// grpc resolver和balancer部分的源码，编程思想：
//  面向接口编程
//  代理模式（wrapper）
//  工厂模式

const schema = "customGrpcResolver"

var cli *clientv3.Client

// 实现了 Builder 接口 & Resolver 接口
type etcdResolver struct {
	// 注册中心地址
	rawAddr string
	// 服务连接信息封装，包含可用地址列表、服务状态等信息。
	// 当有连接信息的时候，通过回调持有grpc ClientConn ,完成连接通知 ,todo
	cc resolver.ClientConn
}

// NewResolver initialize an etcd client
func NewResolver(etcdAddr string) resolver.Builder {

}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

}

func (r etcdResolver) Scheme() string {
	return schema
}

func (r etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	log.Println("ResolveNow !!!  ") // TODO check
}

func (r etcdResolver) Close() {
	log.Println("close !!!")
}

func (r *etcdResolver) watch(keyPrefix string) {
	var addrList []resolver.Address

	getResp, err := cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	} else {
		for i := range getResp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(getResp.Kvs[i].Key), keyPrefix)})
		}
	}

	r.cc.NewAddress(addrList)
	// r.cc.UpdateState(resolver.State{addrList,nil,nil})
	rch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					// 服务上线，服务地址追加
					addrList = append(addrList, resolver.Address{Addr: addr})
					// NewAddress is called by resolver to notify ClientConn a new list
					// of resolved addresses.
					r.cc.NewAddress(addrList)
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.cc.NewAddress(addrList)
				}
			}
			//log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
