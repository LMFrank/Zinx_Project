package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/zinx/utils"
	"zinx/zinx/ziface"
)

/*
	Server服务器模块
*/
type Server struct {
	Name       string              // 服务器名称
	IPVersion  string              // 服务器绑定的 ip 版本
	IP         string              // 服务器监听的 ip 地址
	Port       int                 // 服务器监听的端口
	msgHandler ziface.IMsgHandler  // 当前 Server 的消息管理模块，用来绑定 MsgId 和对应的处理方法
	ConnMgr    ziface.IConnManager // 当前 Server 的连接管理器
}

func NewServer() ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listener at IP: %s, Port: %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	// fmt.Printf("[Start] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 开启消息队列及 Worker 工作池
		s.msgHandler.StartWorkerPool()

		// 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error:", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "error:", err)
			return
		}
		fmt.Printf("Start Zinx server: [%s] success, Listening...\n", s.Name)

		var cid uint32 = 0

		// 3. 阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			// 判断当前连接数量是否超过最大个数
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too many connections!")
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和 conn 进行绑定，得到连接模块
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("[Zinx] STOP server:", s.Name)
	s.ConnMgr.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}
