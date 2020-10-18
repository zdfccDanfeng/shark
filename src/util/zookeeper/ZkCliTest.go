package zookeeper

// 测试zk连接
import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

var (
	hosts       = []string{"127.0.0.1:2181"}
	path        = "/wtzk"
	flags int32 = zk.FlagEphemeral
	data        = []byte("zk data 001")
	acls        = zk.WorldACL(zk.PermAll)
)

func TestConn() {
	// 创建zk连接地址
	hosts := []string{"127.0.0.1:2181"}
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	println(conn.Server())
}

// 增
func add(conn *zk.Conn) {
	var data = []byte("test value")
	// flags有4种取值：
	// 0:永久，除非手动删除
	// zk.FlagEphemeral = 1:短暂，session断开则该节点也被删除
	// zk.FlagSequence  = 2:会自动在节点后面添加序号
	// 3:Ephemeral和Sequence，即，短暂且自动添加序号
	var flags int32 = 0
	// 获取访问控制权限
	acls := zk.WorldACL(zk.PermAll)
	s, err := conn.Create(path, data, flags, acls)
	if err != nil {
		fmt.Printf("创建失败: %v\n", err)
		return
	}
	fmt.Printf("创建: %s 成功", s)
}

// 查
func get(conn *zk.Conn) {
	data, _, err := conn.Get(path)
	if err != nil {
		fmt.Printf("查询%s失败, err: %v\n", path, err)
		return
	}
	fmt.Printf("%s 的值为 %s\n", path, string(data))
}

// 删改与增不同在于其函数中的version参数,其中version是用于 CAS支持
// 可以通过此种方式保证原子性
// 改
func modify(conn *zk.Conn) {
	new_data := []byte("hello zookeeper")
	_, sate, _ := conn.Get(path)
	_, err := conn.Set(path, new_data, sate.Version)
	if err != nil {
		fmt.Printf("数据修改失败: %v\n", err)
		return
	}
	fmt.Println("数据修改成功")
}

// 删
func del(conn *zk.Conn) {
	_, sate, _ := conn.Get(path)
	err := conn.Delete(path, sate.Version)
	if err != nil {
		fmt.Printf("数据删除失败: %v\n", err)
		return
	}
	fmt.Println("数据删除成功")
}

// 测试zk增删改查功能。。
func TestCrudZk() {
	// 创建zk连接地址
	hosts := []string{"127.0.0.1:2181"}
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	/* 增删改查 */
	//add(conn)
	//get(conn)
	//modify(conn)
	del(conn)
	get(conn)
}

// watch机制#
//Java API中是通过Watcher实现的，在go-zookeeper中则是通过Event。道理都是一样的
//全局监听：
//调用zk.WithEventCallback(callback)设置回调

func TestWatch() {
	// 创建监听的option，用于初始化zk
	eventCallbackOption := zk.WithEventCallback(callback)
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5, eventCallbackOption)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 开始监听path
	_, _, _, err = conn.ExistsW(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 触发创建数据操作
	create(conn, path, data)

	//再次监听path
	_, _, _, err = conn.ExistsW(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 触发删除数据操作
	del2(conn, path)
}

// zk watch 回调函数
func callback(event zk.Event) {
	// zk.EventNodeCreated
	// zk.EventNodeDeleted
	fmt.Println("###########################")
	fmt.Println("path: ", event.Path)
	fmt.Println("type: ", event.Type.String())
	fmt.Println("state: ", event.State.String())
	fmt.Println("---------------------------")
}

// 创建数据
func create(conn *zk.Conn, path string, data []byte) {
	_, err := conn.Create(path, data, flags, acls)
	if err != nil {
		fmt.Printf("创建数据失败: %v\n", err)
		return
	}
	fmt.Println("创建数据成功")
}

// 删除数据
func del2(conn *zk.Conn, path string) {
	_, stat, _ := conn.Get(path)
	err := conn.Delete(path, stat.Version)
	if err != nil {
		fmt.Printf("删除数据失败: %v\n", err)
		return
	}
	fmt.Println("删除数据成功")
}

// 部分监听：
//1.调用conn.ExistsW(path) 或GetW(path)为对应节点设置监听，该监听只生效一次
//2.开启一个协程处理chanel中传来的event事件
//（注意：watchCreataNode一定要放在一个协程中，不能直接在main中调用，不然会阻塞main）
func TestPartialWatch() {
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 开始监听path
	_, _, event, err := conn.ExistsW(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 协程调用监听事件
	go watchZkEvent(event)

	// 触发创建数据操作
	create2(conn, path, data)
}

// zk 回调函数
func watchZkEvent(e <-chan zk.Event) {
	event := <-e
	fmt.Println("###########################")
	fmt.Println("path: ", event.Path)
	fmt.Println("type: ", event.Type.String())
	fmt.Println("state: ", event.State.String())
	fmt.Println("---------------------------")
}

// 创建数据
func create2(conn *zk.Conn, path string, data []byte) {
	_, err := conn.Create(path, data, flags, acls)
	if err != nil {
		fmt.Printf("创建数据失败: %v\n", err)
		return
	}
	fmt.Println("创建数据成功")
}

//客户端随机hostname支持#
//最终就是Round Robin策略
func TestRobinLb() {
	var hosts = []string{"host01:2181", "host02:2181", "host03:2181"}
	hostPro := new(zk.DNSHostProvider)
	//先初始化
	err := hostPro.Init(hosts)

	if err != nil {
		fmt.Println(err)
		return
	}
	//获得host
	server, retryStart := hostPro.Next()
	fmt.Println(server, retryStart)
	//连接成功后会调用
	hostPro.Connected()
}
