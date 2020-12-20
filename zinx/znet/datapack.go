package znet

import (
	"bytes"
	"encoding/binary"
	"zinx/zinx/ziface"
)

// 封包拆包类实例，暂时不需要成员
// TLV 模型：head(dataLen + id) + body(data)
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// dataLen uint32（4 字节） + id uint32（4 字节）
	return 8
}

//封包方法(压缩数据)
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放 byte 字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 写 dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写 msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写 data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法(解压数据)
// 只需要读出 Head 信息，再根据 Head 信息里的 data 长度读取
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入读取二进制数据的 ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压 head 信息，得到 dataLen 和 msgID
	msg := &Message{}

	// 读 dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读 msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断 dataLen 的长度是否超出我们允许的最大包长度
	//if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
	//	return nil, errors.New("too large msg data")
	//}

	return msg, nil
}
