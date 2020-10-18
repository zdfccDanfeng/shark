package client

import (
	"context"
	"fmt"
	"github.com/shark/src/util/zookeeper"
	pb "github.com/shark/src/util/zookeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"time"
)

// 参考：@see https://blog.csdn.net/Edu_enth/article/details/104016307

func StartClient() {
	r := zookeeper.NewResolver("localhost:2378") // etcd 地址
	resolver.Register(r)
	// 如果请求的负载均衡地址，会基于grpclb返回一个有效的连接地址：@see https://colobu.com/2017/03/25/grpc-naming-and-load-balance/
	conn, err := grpc.Dial(r.Scheme()+"://author/project/test", grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	for {
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "hello"}, grpc.WaitForReady(true))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp)
		}

		<-time.After(time.Second)
	}
}
