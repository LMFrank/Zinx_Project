package utils

import (
	"encoding/json"
	"errors"
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
	Version          string // 当前 Zinx 版本号
	MaxConn          int    // 当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 // 工作池的 Worker 数量
	MaxPackageSize   uint32 // 数据包的最大值
	MaxWorkerTaskLen uint32 // Worker 对应负责的任务队列最大任务存储数量
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() error {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		return errors.New("can't reload zinx.json")
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		return errors.New("can't unmarshal zinx.json")
	}
	return nil
}

func init() {
	// 默认配置
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.9",
		TcpPort:          9000,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		WorkerPoolSize:   10,
		MaxPackageSize:   4096,
		MaxWorkerTaskLen: 1024,
	}

	// 从配置文件中加载用户自定义的配置参数
	err := GlobalObject.Reload()
	if err != nil {
		return
	}
}
