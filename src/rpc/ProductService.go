package rpc

import (
	"context"
	"github.com/shark/src/rpc/proto/product"
	"log"
)

// productService实现类

type ProductRrcSerciceImpl struct {
	ProductSercice.ProductServiceService
}

// 服务端实现，实现rpc接口
func (this *ProductRrcSerciceImpl) QueryProdInfoDetail(context context.Context, req *ProductSercice.ProdctInfo) (*ProductSercice.Response, error) {
	log.Printf("productInfo is : %v\n", req)
	return &ProductSercice.Response{Ok: 12}, nil
}
