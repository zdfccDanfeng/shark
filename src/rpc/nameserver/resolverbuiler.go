package nameserver

import (
	"context"
	"google.golang.org/grpc/resolver"
)

// 自定义名称服务

type MyCustomNsBuilder struct{}

func NewBuilder() resolver.Builder {
	return &MyCustomNsBuilder{}
}

// 采用了builer设计模式，表示和构建分离
func (this *MyCustomNsBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	ctx, cancel := context.WithCancel(context.Background())

}

func (this *MyCustomNsBuilder) Scheme() string {

}
