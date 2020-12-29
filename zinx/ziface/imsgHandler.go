package ziface

type IMsgHandler interface {
	// 调度/执行对应的 Router 消息处理方法
	DoMsgHandler(request IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)

	// 启动 worker 工作池
	StartWorkerPool()

	// 将消息交给 TaskQueue，由 Worker 进行处理
	SendMsgToTaskQueue(request IRequest)
}
