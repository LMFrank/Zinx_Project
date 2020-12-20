package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	MockServer()
	MockClient()
}

/*
	模拟的 server 端
*/
func MockServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	// 创建服务器 goroutine，负责从客户端 goroutine 读取粘包的数据，然后进行解析
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Server accept err:", err)
			}

			// 处理客户端请求
			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					// 先读出流中的 head 部分
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData) // ReadFull 会把msg填充满为止
					if err != nil {
						fmt.Println("Read head error")
					}

					// 将 headData 字节流拆包到 msg 中
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("Server unpack err:", err)
						return
					}
					if msgHead.GetDataLen() > 0 {
						// msg 是有 data 数据的，需要再次读取 data 数据
						if msg, ok := msgHead.(*Message); ok {
							msg.Data = make([]byte, msg.GetDataLen())
							// 根据 dataLen 从 io 中读取字节流
							_, err := io.ReadFull(conn, msg.Data)
							if err != nil {
								fmt.Println("Server unpack data err:", err)
								return
							}

							fmt.Println("==> Recv Msg: ID =", msg.Id, ", len =", msg.DataLen, ", data =", string(msg.Data))
						}
					}
				}
			}(conn)
		}
	}()
}

/*
	模拟的 client 端
*/
func MockClient() {
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:8000")
		if err != nil {
			fmt.Println("client dial err:", err)
			return
		}

		dp := NewDataPack()

		// 封装一个msg1包
		msg1 := &Message{
			Id:      1,
			DataLen: 5,
			Data:    []byte{'h', 'e', 'l', 'l', 'o'},
		}

		sendData1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("Client pack msg1 err:", err)
			return
		}

		msg2 := &Message{
			Id:      1,
			DataLen: 7,
			Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
		}
		sendData2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("Client pack msg2 err:", err)
			return
		}

		//将sendData1，和 sendData2 拼接一起，组成粘包
		sendData1 = append(sendData1, sendData2...)

		//向服务器端写数据
		_, err = conn.Write(sendData1)
		if err != nil {
			panic(err)
		}
	}()

	//客户端阻塞
	select {
	case <-time.After(time.Second):
		return
	}
}
