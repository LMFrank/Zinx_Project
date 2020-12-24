package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/zinx/ziface"
)

// 连接模块
type Connection struct {
	Conn       *net.TCPConn       // 当前连接的 socket TCP 套接字
	ConnID     uint32             // 当前连接的 ID
	isClosed   bool               // 当前连接的状态
	ExitChan   chan bool          // 告知当前连接已经退出/停止的 channel
	MsgHandler ziface.IMsgHandler // 当前 server 的消息管理模块，用来绑定 MsgID 和对应的处理业务 API 关系
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
	}
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID =", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()

		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("Read msg head error:", err)
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack error:", err)
			break
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("Read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前 conn 数据的 Request 请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 从路由中找到注册绑定的 Conn 对应的 router 调用
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start()... ConnID =", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
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

// 将发送给 client 端的数据先进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id =", msgId)
		return errors.New("pack error msg")
	}

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id:", msgId, "error:", err)
		return errors.New("conn write error")
	}

	return nil
}
