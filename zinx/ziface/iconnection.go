package ziface

import "net"

// 连接
type IConnection interface {
	// 启动连接 让当前的连接准备开始工作
	Start()

	// 停止连接 结束当前连接的工作
	Stop()

	// 获取当前连接绑定的 socket
	GetTCPConnection() *net.TCPConn

	// 获取当前连接模块的 ID
	GetConnID() uint32

	// 获取远程客户端的 TCP 状态 IP Port
	RemoteAddr() net.Addr

	//直接将Message数据发送数据给远程的TCP客户端(无缓冲)
	SendMsg(msgId uint32, data []byte) error
}

// 处理连接业务
type HandleFunc func(*net.TCPConn, []byte, int) error
