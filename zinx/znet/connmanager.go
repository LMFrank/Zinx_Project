package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/zinx/ziface"
)

/*
	连接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接信息
	connLock    sync.RWMutex                  // 读写连接的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn

	fmt.Println("Connection ADD ConnID =", conn.GetConnID(), "successfully: conn num =", c.Len())
}

// 删除连接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, conn.GetConnID())

	fmt.Println("Connection REMOVE ConnID =", conn.GetConnID(), "successfully: conn num =", c.Len())
}

// 根据 ConnID 获取连接
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 获取当前连接数量
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// 清除并停止所有连接
func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(c.connections, connID)
	}

	fmt.Println("CLEAR All Connections successfully: conn num =", c.Len())
}
