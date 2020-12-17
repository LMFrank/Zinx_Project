package znet

import (
	"fmt"
	"net"
	"zinx/zinx/ziface"
)

// Server服务器模块
type Server struct {
	Name      string // 服务器名称
	IPVersion string // 服务器绑定的ip版本
	IP        string // 服务器监听的ip地址
	Port      int    // 服务器监听的端口
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err", err)
			return
		}
		fmt.Println("start Zinx server success,", s.Name, "success, Listening...")

		// 3. 阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			// 做一个基本的回显业务用于测试
			go func() {
				for {
					buf := make([]byte, 512)
					count, err := conn.Read(buf)
					if err != nil {
						fmt.Println("receive buf error ", err)
						continue
					}
					fmt.Printf("receive client buf %s, count %d\n", buf, count)

					if _, err := conn.Write(buf[0:count]); err != nil {
						fmt.Println("write back buf error ", err)
						continue
					}

				}
			}()
		}
	}()
}

func (s *Server) Stop() {
	panic("implement me")
}

func (s *Server) Serve() {
	s.Start()
	select {}
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      9000,
	}
}
