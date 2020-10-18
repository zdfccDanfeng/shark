package zookeeper

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"time"
)

// #简单的分布式server 目前分布式系统已经很流行了，一些开源框架也被广泛应用，如dubbo、Motan等。对于一个分布式服务，最基本的一项功能就是服务的注册和发现，
//而利用zk的EPHEMERAL节点则可以很方便的实现该功能。EPHEMERAL节点正如其名，是临时性的，其生命周期是和客户端会话绑定的，当会话连接断开时，节点也会被删除。下边我们就来实现一个简单的分布式server：
//server：
//服务启动时，创建zk连接，并在go_servers节点下创建一个新节点，节点名为"ip:port"，完成服务注册
//服务结束时，由于连接断开，创建的节点会被删除，这样client就不会连到该节点
//client：
//
//先从zk获取go_servers节点下所有子节点，这样就拿到了所有注册的server
//从server列表中选中一个节点（这里只是随机选取，实际服务一般会提供多种策略），创建连接进行通信

// 为了演示，我们每次client连接server，获取server发送的时间后就断开。主要代码如下： zkutil:用来处理zk操作
func GetConnect() (conn *zk.Conn, err error) {
	conn, _, err = zk.Connect(hosts, time.Second)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func RegistServer(conn *zk.Conn, host string) (err error) {
	_, err = conn.Create("/go_servers/"+host, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return
}

func GetServerList(conn *zk.Conn) (list []string, err error) {
	list, _, err = conn.Children("/go_servers")
	return
}

// server: server端操作
func mockStartServer() {
	go starServer("127.0.0.1:8897")
	go starServer("127.0.0.1:8898")
	go starServer("127.0.0.1:8899")

	a := make(chan bool, 1)
	<-a
}

func starServer(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
	fmt.Println(tcpAddr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//注册zk节点q
	conn, err := GetConnect()
	if err != nil {
		fmt.Printf(" connect zk error: %s ", err)
	}
	defer conn.Close()
	err = RegistServer(conn, port)
	if err != nil {
		fmt.Printf(" regist node error: %s ", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
			continue
		}
		go handleCient(conn, port)
	}

	fmt.Println("aaaaaa")
}

func handleCient(conn net.Conn, port string) {
	defer conn.Close()

	daytime := time.Now().String()
	conn.Write([]byte(port + ": " + daytime))
}

// ============ client: 客户端操作

func mockStartClient() {
	for i := 0; i < 100; i++ {
		startClient()

		time.Sleep(1 * time.Second)
	}
}

func startClient() {
	// service := "127.0.0.1:8899"
	//获取地址
	serverHost, err := getServerHost()
	if err != nil {
		fmt.Printf("get server host fail: %s \n", err)
		return
	}

	fmt.Println("connect host: " + serverHost)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverHost)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()

	_, err = conn.Write([]byte("timestamp"))
	checkError(err)

	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println(string(result))

	return
}

func checkError(err error) {

}

func getServerHost() (host string, err error) {
	conn, err := GetConnect()
	if err != nil {
		fmt.Printf(" connect zk error: %s \n ", err)
		return
	}
	defer conn.Close()
	serverList, err := GetServerList(conn)
	if err != nil {
		fmt.Printf(" get server list error: %s \n", err)
		return
	}

	count := len(serverList)
	if count == 0 {
		err = errors.New("server list is empty \n")
		return
	}

	//随机选中一个返回
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	host = serverList[r.Intn(3)]
	return
}
