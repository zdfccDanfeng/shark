package rpc

import (
	"context"
	"fmt"
	"github.com/shark/src/rpc/proto/product"
	"google.golang.org/grpc"
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
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("resp is :", resp)
}
