package ziface

// 路由抽象接口
// 路由里的数据均为 IRequest
type IRouter interface {
	// 处理 conn 业务之前的钩子方法 Hook
	PreHandle(request IRequest)
	// 处理 conn 业务的主方法 Hook
	Handle(request IRequest)
	// 处理 conn 业务之后的钩子方法 Hook
	PostHandle(request IRequest)
}
