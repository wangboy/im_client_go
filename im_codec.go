package main

import (
	"awesomeProject/common"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	"reflect"
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
	protocol := msgObj.(common.Protocol)

	ret := common.Octets{}
	protocol.Marshal(&ret)

	return ret.Data, nil
	//return json.Marshal(msgObj.(common.Protocol))

}

// 将JSON的字节数组解码为结构体
func (self *ImCodec) Decode(data interface{}, msgObj interface{}) error {

	codecLogger.Infoln(" rece data ", data, msgObj)

	bytes := data.([]byte)

	b := common.Octets{bytes, 0, 0}
	msgId := b.ReadInt()

	msg := msgs[msgId]
	protocol := msg.(common.Protocol)
	protocol.Unmarshal(&b)

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

//
//type common.Protocol interface {
//	getMsgId() int32
//	marshal(out *common.Octets)
//	unmarshal(in *common.Octets)
//}

var (
	msgs = make(map[int32]reflect.Type)
)
