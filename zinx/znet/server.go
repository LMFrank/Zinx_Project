package znet

import (
	"fmt"
	"net"
	"zinx/zinx/ziface"
)

// Server服务器模块
type Server struct {
	Name      string         // 服务器名称
	IPVersion string         // 服务器绑定的 ip 版本
	IP        string         // 服务器监听的 ip 地址
	Port      int            // 服务器监听的端口
	Router    ziface.IRouter // 当前的 server 添加一个router，server 注册的连接对应的处理业务
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " error:", err)
			return
		}
		fmt.Println("start Zinx server success,", s.Name, "success, Listening...")
		var cid uint32
		cid = 0

		// 3. 阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			// 将处理新连接的业务方法和 conn 进行绑定，得到连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	panic("implement me")
}

// 运行服务器
func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success")
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      9000,
		Router:    nil,
	}
}
