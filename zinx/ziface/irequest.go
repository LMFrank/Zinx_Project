package ziface

// IRequest 接口：
// 把客户端请求的连接信息和请求的数据包封装到一个 Request 中
type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection

	// 得到请求的消息数据
	GetData() []byte
}
