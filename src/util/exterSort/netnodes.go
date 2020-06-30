package exterSort

import (
	"bufio"
	"fmt"
	"net"
)

// 从网络节点接受数据写入channel
func NetWorkSink(addr string, in <-chan int) {
	// 开启监听
	listener, e := net.Listen("tcp", addr)
	if e != nil {
		panic(e)
	}
	// 后端的goroutine 进行处理
	go func() {
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("conn is : remote addr :", conn.RemoteAddr(), " local is :", conn.LocalAddr())
		defer conn.Close()
		writer := bufio.NewWriter(conn)
		defer writer.Flush()
		Sink(writer, in)
	}()
}
func NetWorkSource(addr string) <-chan int {
	out := make(chan int)
	go func() {
		conn, e := net.Dial("tcp", addr)
		if e != nil {
			panic(e)
		}
		source := ReadSource(bufio.NewReader(conn), -1)
		for v := range source {
			out <- v
		}
		close(out)
	}()
	return out
}
