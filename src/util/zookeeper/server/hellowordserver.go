package server

import (
	"context"
	"github.com/shark/src/util/zookeeper"
	pb "github.com/shark/src/util/zookeeper/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const svcName = "project/test"

var addr = "127.0.0.1:50051"

func StartServer() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	defer s.GracefulStop()

	pb.RegisterGreeterServer(s, helloRpcImpl{})
	// 注册服务信息到zk / etcd ...
	go zookeeper.Register("127.0.0.1:2379", svcName, addr, 5)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		zookeeper.UnRegister(svcName, addr)
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}

type helloRpcImpl struct {
}

func (h helloRpcImpl) SayHello(context context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// req.Data = req.Data + ", from:" + addr
	return &pb.HelloReply{Message: addr}, nil
}
