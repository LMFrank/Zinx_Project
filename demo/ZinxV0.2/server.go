package main

import "zinx/zinx/znet"

/*
	模拟服务端
*/
func main() {
	s := znet.NewServer("[Zinx V0.2]")
	s.Serve()
}
