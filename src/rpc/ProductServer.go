package rpc

import (
	"context"
	"fmt"
	"github.com/shark/src/rpc/proto/product"
	"google.golang.org/grpc"
	"log"
	"net"
)

// 实现RPC接口。。。
func QueryProdInfoDetail(context context.Context, req *ProductSercice.ProdctInfo) (*ProductSercice.Response, error) {
	log.Printf("productInfo is : %v\n", req)
	return &ProductSercice.Response{Ok: 12}, nil
}
func QueryBatchProdInfoDetail(req *ProductSercice.ProdctInfo, stream ProductSercice.ProductService_QueryBatchProdInfoDetailServer) error {
	// 每两条发送一次。。

	for i := 0; i < 100; i++ {
		if i%2 == 0 && i > 0 {
			err := stream.Send(&ProductSercice.Response{Desc: "xx"})
			if err != nil {
				// 单次发送处理出错，。。
				return err
			}
		}
	}
	return nil
}

func NewProductServer() {
	fmt.Println("server begin .. start !!!!")
	rpcServer := grpc.NewServer()
	// 注册服务
	ProductSercice.RegisterProductServiceService(rpcServer, &ProductSercice.ProductServiceService{QueryProdInfoDetail: QueryProdInfoDetail})
	lis, _ := net.Listen("tcp", ":8082")
	rpcServer.Serve(lis)
	fmt.Println("server started !!!!")
}
