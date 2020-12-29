package znet

import (
	"fmt"
	"strconv"
	"zinx/zinx/utils"
	"zinx/zinx/ziface"
)

/*
	消息处理模块
*/
type MsgHandler struct {
	APIs           map[uint32]ziface.IRouter // 存放每个 MsgID 对应的处理方法
	TaskQueue      []chan ziface.IRequest    // 负责 worker 取任务的消息队列
	WorkerPoolSize uint32                    // 业务工作 worker 池的 worker 数量
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// 调度/执行对应的 Router 消息处理方法
func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("API msgId =", request.GetMsgID(), "is not found!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前 msg 绑定的 API 处理方法是否存在
	if _, ok := m.APIs[msgID]; ok {
		panic("Repeated API , msgId = " + strconv.Itoa(int(msgID)))
	}

	// 添加 msg 与 API 的绑定关系
	m.APIs[msgID] = router
	fmt.Println("Add API MsgID =", msgID, "success!")
}

// 启动一个 Worker 工作池
func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 给当前 worker 对应的任务队列开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前 Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个 Worker 工作流程
func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID =", workerID, "is started.")
	// 不断的阻塞等待队列中的消息
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 将消息交给 TaskQueue，由 Worker 进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 采用平均分配消息给 Worker
	// 根据客户端建立的 connID 来进行分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID =", request.GetConnection().GetConnID(), "request msgID=", request.GetMsgID(), "to workerID=", workerID)
	// 将请求消息发送给任务队列
	m.TaskQueue[workerID] <- request
}
