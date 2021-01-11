package contextusage

// Golang 使用树形派生的方式构造 Context，通过在不同过程 [1] 中传递 deadline 和 cancel 信号，来管理处理某个任务所涉及到的一组 goroutine 的生命周期，防止 goroutine 泄露。
//并且可以通过附加在 Context 上的 Value 来传递/共享一些跨越整个请求间的数据。

// Context 最常用来追踪 RPC/HTTP 等耗时的、跨进程的 IO 请求的生命周期，从而让外层调用者可以主动地或者自动地取消该请求，进而告诉子过程回收用到的所有 goroutine 和相关资源。
// Context 本质上是一种在 API 间树形嵌套调用时传递信号的机制。本文将从接口、派生、源码分析、使用等几个方面来逐一解析 Context。

// // Context 用以在多 API 间传递 deadline、cancelation 信号和请求的键值对。
//// Context 中的方法能够安全的被多个 goroutine 并发调用。

//type Context interface {
//	// Done 返回一个只读 channel，该 channel 在 Context 被取消或者超时时关闭
//	Done() <-chan struct{}
//
//	// Err 返回 Context 结束时的出错信息
//	Err() error
//
//	// 如果 Context 被设置了超时，Deadline 将会返回超时时限。
//	Deadline() (deadline time.Time, ok bool)
//
//	// Value 返回关联到相关 Key 上的值，或者 nil.
//	Value(key interface{}) interface{}
//}
// Done() 方法返回一个只读的 channel，当 Context 被主动取消或者超时自动取消时，该 Context 所有派生 Context 的 done channel 都被 close 。所有子过程通过该字段收到 close 信号后，
// 应该立即中断执行、释放资源然后返回
// Value() 返回绑定在该 Context 链（我称为回溯链，下文会展开说明）上的给定的 Key 的值，如果没有，则返回 nil。注意，不要用于在函数中传参，其本意在于共享一些横跨整个 Context 生命周期范围的值。
//Key 可以是任何可比较的类型。为了防止 Key 冲突，最好将 Key 的类型定义为非导出类型，然后为其定义访问器

import "context"

// User 是要存于 Context 中的 Value 类型.
type User struct{}

// key 定义为了非导出类型，以避免和其他 package 中的 key 冲突
type key int

// userKey 是 Context 中用来关联 user.User 的 key，是非导出变量
// 客户端需要用 user.NewContext 和 user.FromContext 构建包含
// user 的 Context 和从 Context 中提取相应 user
var userKey key

// NewContext 返回一个带有用户值 u 的 Context.
func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext 从 Context 中提取 user，如果有的话.
func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

// Context 派生
// Context 设计之妙在于可以从已有 Context 进行树形派生，以管理一组过程的生命周期。我们上面说了单个 Context 实例是不可变的，但可以通过 context 包提供的三种方法：WithCancel 、
// WithTimeout 和 WithValue 来进行派生并附加一些属性（可取消、时限、键值），以构造一组树形组织的 Context。
// 当根 Context 结束时，所有由其派生出的 Context 也会被一并取消。也就是说，父 Context 的生命周期涵盖所有子 Context 的生命周期。
// context.Background() 通常用作根节点，它不会超时，不能被取消
// 通过 WithCancel 从 context.Background() 派生出的 Context 要注意在对应过程完结后及时 cancel，否则会造成 Context 泄露。
// context.Background() 和 context.TODO() 返回的都是 emptyCtx 的实例。但其语义略有不同。前者做为 Context 树的根节点，后者通常在不知道用啥时用
// Background 返回一个空 Context。它不能被取消，没有时限，没有附加键值。Background 通常用在
// main函数、init 函数、test 入口，作为某个耗时过程的根 Context。
// @see https://zhuanlan.zhihu.com/p/163684835
