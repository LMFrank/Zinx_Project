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

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
