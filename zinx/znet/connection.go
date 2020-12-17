package znet

import (
	"fmt"
	"net"
	"zinx/zinx/ziface"
)

// 连接模块
type Connection struct {
	Conn      *net.TCPConn      // 当前连接的 socket TCP 套接字
	ConnID    uint32            // 当前连接的ID
	isClosed  bool              // 当前连接的状态
	handleAPI ziface.HandleFunc // 当前连接所绑定的处理业务方法 API
	ExitChan  chan bool         // 告知当前连接已经退出/停止的 channel
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, " Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端数据到 buf 中
		buf := make([]byte, 512)
		count, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Receive buf error:", err)
			continue
		}

		// 调用当前连接所绑定的 HandleAPI
		if err := c.handleAPI(c.Conn, buf, count); err != nil {
			fmt.Println("ConnID", c.ConnID, " handle is error:", err)
			break
		}
	}
}

// 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start()... ConnID=", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()

	// 启动从当前连接写数据的业务
}

// 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭 socket 连接
	c.Conn.Close()
	// 回收资源
	close(c.ExitChan)
}

// 获取当前连接绑定的 socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的 ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 TCP 状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 将数据发送给远程客户端
func (c *Connection) Send(data []byte) error {
	panic("implement me")
}

func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		handleAPI: callbackApi,
		ExitChan:  make(chan bool, 1),
	}
}
