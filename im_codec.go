package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
)

var codecLogger = golog.New("codec")

type ImCodec struct {
}

// 编码器的名称
func (self *ImCodec) Name() string {
	return "ImCodes"
}

func (self *ImCodec) MimeType() string {
	return "application/ImCodes"
}

// 将结构体编码为JSON的字节数组
func (self *ImCodec) Encode(msgObj interface{}, ctx cellnet.ContextSet) (data interface{}, err error) {
	protocol := msgObj.(Protocol)

	ret := ImBuffer{}
	protocol.marshal(&ret)

	return ret.data, nil
	//return json.Marshal(msgObj.(Protocol))

}

// 将JSON的字节数组解码为结构体
func (self *ImCodec) Decode(data interface{}, msgObj interface{}) error {

	codecLogger.Infoln(" rece data ", data, msgObj)

	bytes := data.([]byte)

	b := ImBuffer{bytes, 0, 0}
	msgId := b.readInt()

	msg := msgs[msgId]
	protocol := msg.(Protocol)
	protocol.unmarshal(&b)

	return nil
}
//
//func init() {
//	imCodec := new(ImCodec)
//	fmt.Println(" in in init codec ", imCodec)
//
//	// 注册编码器
//	codec.RegisterCodec(imCodec)
//}
