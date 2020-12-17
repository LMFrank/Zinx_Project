package main

import (
	"fmt"
	"zinx/zinx/ziface"
	"zinx/zinx/znet"
)

/*
	模拟服务端
*/

// ping test
type PingRouter struct {
	znet.BaseRouter
}

// Test PreRouter
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("Call back ping ping ping error")
	}
}

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("Call back ping ping ping error")
	}
}

// Test PostHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("Call back ping ping ping error")
	}
}

func main() {
	s := znet.NewServer("[Zinx V0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
