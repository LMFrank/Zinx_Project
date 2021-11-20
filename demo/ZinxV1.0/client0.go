package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/zinx/znet"
)

/*
	模拟客户端
*/
func main() {
	fmt.Println("Client0 Test... start")
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8500")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 发封包 message 消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx_V1.0 client0 Test Message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error:", err)
			return
		}

		// 先读出流中的 head 部分
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err = io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("Read head error")
			break
		}

		// 将 headData 字节流拆包到 msg 中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("Server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// msg 是有 data 数据的，需要再次读取 data 数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据 dataLen 从 io 中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("Server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Server Msg: ID =", msg.Id, ", len =", msg.DataLen, ", data =", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
