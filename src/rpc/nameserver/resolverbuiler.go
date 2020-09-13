package nameserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc/resolver"
)

// 自定义名称服务

type MyCustomNsBuilder struct{}

// 特殊的函数init函数，先于main函数执行，实现包级别的一些初始化操作
// init 将定义好的 NS Builder 注册到 resolver 包中
// 初始化不能采用初始化表达式初始化的变量。
// 程序运行前的注册。
// 实现sync.Once功能。
// init函数先于main函数自动执行，不能被其他函数调用；
// 每个包可以有多个init函数；
// 采用 Builder 模式在包初始化时创建并注册构造 nsResover 的 MyCustomResolver 实例 --> todo [业务代码使用的时候引入该包名，会在包初始化的时候，将该自定义的builer注册到resolver的全局m映射里面]
// 当客户端通过 Dial 方法对指定服务进行拨号时，grpc resolver 查找注册的 Builder 实例调用其 Build() 方法构建自定义 nsResolver。
func init() {
	resolver.Register(NewBuilder())
}

func NewBuilder() resolver.Builder {
	return &MyCustomNsBuilder{}
}

// URI 返回某个服务的统一资源描述符（URI），这个 URI 可以从 nsResolver 中查询实例列表
// URI 设计时可以遵循 RFC-3986(https://tools.ietf.org/html/rfc3986) 规范，
// 比如本例中 ns 格式为：ns://callerService:@calleeService
// 其中 ns 为协议名，callerService 为订阅方服务名（即主调方服务名），calleeService 为发布方服务名（即被调方服务名）
func URI(callerService, calleeService string) string {
	return fmt.Sprintf("ns://%s:@%s", callerService, calleeService)
}

// 采用了builer设计模式，表示和构建分离
func (this *MyCustomNsBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	ctx, cancel := context.WithCancel(context.Background())

	r := &MyCustomResolver{
		target: target,
		cc:     cc,
		ctx:    ctx,
		cancel: cancel,
	}
	// 启动协程，响应指定 Name 服务实例变化
	go r.watcher()
	return r, nil

}

// Scheme 实现了 resolver.Builder.Scheme 方法
// Scheme 方法定义了 ns resolver 的协议名
func (this *MyCustomNsBuilder) Scheme() string {
	return "my_custom_ns"
}
