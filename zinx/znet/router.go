package znet

import "zinx/zinx/ziface"

// 路由基类，根据需要对基类方法重写
type BaseRouter struct{}

// 处理 conn 业务之前的钩子方法 Hook
func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

// 处理 conn 业务的主方法 Hook
func (b *BaseRouter) Handle(request ziface.IRequest) {}

// 处理 conn 业务之后的钩子方法 Hook
func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
