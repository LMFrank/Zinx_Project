package znet

import (
	"fmt"
	"strconv"
	"zinx/zinx/ziface"
)

/*
	消息处理模块
*/
type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个 MsgID 对应的处理方法
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 调度/执行对应的 Router 消息处理方法
func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("Api msgId =", request.GetMsgID(), "is not found!")
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
	if _, ok := m.Apis[msgID]; ok {
		panic("Repeated api , msgId = " + strconv.Itoa(int(msgID)))
	}

	// 添加 msg 与 API 的绑定关系
	m.Apis[msgID] = router
	fmt.Println("Add api MsgID =", msgID, "success!")
}
