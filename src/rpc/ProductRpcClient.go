package rpc

import (
	"context"
	"fmt"
	"github.com/shark/src/rpc/proto/product"
	"google.golang.org/grpc"
	"io"
	"log"
)

func NewProductClient() {
	client, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	procuctServiceClient := ProductSercice.NewProductServiceClient(client)
	resp, err := procuctServiceClient.QueryProdInfoDetail(context.Background(), &ProductSercice.ProdctInfo{Desc: "xx"})
	detail, err := procuctServiceClient.QueryBatchProdInfoDetail(context.Background(), &ProductSercice.ProdctInfo{})
	if err != nil {
		return
	}
	for {
		// 循环读取流里面的数据。。
		recv, err := detail.Recv()
		if err != io.EOF {
			break // 数据发送完毕，但是不是异常错误，退出程序，而不是停掉服务。。
		}
		if err != nil {
			log.Fatal("error") // 停掉服务
		}
		fmt.Println("recv ", recv)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("resp is :", resp)
}
