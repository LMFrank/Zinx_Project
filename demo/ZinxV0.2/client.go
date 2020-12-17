package main

import (
	"fmt"
	"net"
	"time"
)

/*
	模拟客户端
*/
func main() {
	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second)

	// 1. 连接远程服务器，得到一个连接
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	// 2. 连接调用 Write 写数据
	for {
		_, err := conn.Write([]byte("Hello Zinx V0.1..."))
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		count, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}

		fmt.Printf(" server call back : %s, count = %d\n", buf, count)

		time.Sleep(1 * time.Second)
	}
}
