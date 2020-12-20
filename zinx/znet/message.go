package znet

type Message struct {
	DataLen uint32 // 消息的长度
	Id      uint32 // 消息的ID
	Data    []byte // 消息的内容
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		DataLen: uint32(len(data)),
		Id:      id,
		Data:    data,
	}
}

// 获取消息数据段长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息数据段长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

// 设计消息ID
func (m *Message) SetMsgId(msgId uint32) {
	m.Id = msgId
}

// 设计消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
