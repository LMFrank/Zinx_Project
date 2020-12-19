package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/zinx/ziface"
)

// 全局参数及用户自定义变量
type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer // 当前 Zinx 的全局 Server 对象
	Host      string         // 当前服务器主机 IP
	TcpPort   int            // 当前服务器主机监听端口号
	Name      string         // 当前服务器名称

	/*
		Zinx
	*/
	Version       string // 当前 Zinx 版本号
	MaxConn       int    // 当前服务器主机允许的最大链接个数
	MaxPacketSize uint32 // 都需数据包的最大值
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 默认配置
	GlobalObject = &GlobalObj{
		Name:          "ZinxServerApp",
		Version:       "V0.4",
		TcpPort:       9000,
		Host:          "0.0.0.0",
		MaxConn:       1000,
		MaxPacketSize: 4096,
	}

	// 从配置文件中加载用户自定义的配置参数
	GlobalObject.Reload()
}
