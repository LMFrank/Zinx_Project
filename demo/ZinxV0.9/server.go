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

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	fmt.Println("Recv Client Msg: msgID =", request.GetMsgID(), " data =", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

// hello test
type HelloRouter struct {
	znet.BaseRouter
}

// Test Handle
func (p *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	fmt.Println("Recv Client Msg: msgID =", request.GetMsgID(), " data =", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ... ")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is Called ... ")
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
