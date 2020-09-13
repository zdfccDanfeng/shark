package nameserver

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// grpc.Dial() 方法进行封装，方便业务方适用。封装后 dial.go 代码如下所示（严格来说 dial.go 不应该放在 ns 包中，本例中这么做只是为简化包布局，方便理解）

// Dial 封装 `grpc.Dial()` 方法以供业务方代码初始化 *grpc.ClientConn。
// 业务方可使用此 Dial 方法基于主调方服务名、被调方服务名等参数构造 *grpc.ClientConn 实例，
// 随后可在业务代码中使用 *grpc.ClientConn 实例构造桩代码中生成的 grpcServiceClient 并发起 RPC 调用。
func Dial(callerService, calleeService string, dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// 根据 callerService 和 calleeService 构造对应的 URI
	URI := URI(callerService, calleeService)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置拨号配置
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
	dialOpts = append(dialOpts, dialOpts...)

	// 调用 grpc.DialContext() 方法拨号
	// 方法使用 callerService 和 calleeService 构造服务 URI，并使用此 URI 作为参数调用 grpc.DialContext() 方法，来构造 *grpc.ClientConn 实例。
	// grpc.DialContext() 方法接收三个参数：ctx、target、opts，
	// todo 就是根据我们自定义的协议名、callerService、CalleeService 生成的 URI，比如本例中 target 参数值为 ns://my-caller-service:@my-callee-service，其中 ns 为协议名。grpc 可通过协议名查表来获取对应的 resolverBuilder。
	// opts：是一个变长参数，表示拨号配置选项。
	conn, err := grpc.DialContext(
		ctx,
		URI,
		opts...,
	)
	if err != nil {
		// logz.Warn("did not connect", logz.Any("target", URI), logz.E(err))
		return nil, err
	}
	return conn, err
}
