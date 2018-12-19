package main

import (
	"fmt"
	"reflect"
)

var (
	msgs = make(map[int32]reflect.Type)
)

type Protocol interface {
	getMsgId() int32
	marshal(out *ImBuffer)
	unmarshal(in *ImBuffer)
}

type CTestMsg struct {
	msgId   int32
	content string
}

func (msg *CTestMsg) getMsgId() int32 {
	return msg.msgId
}

func (msg *CTestMsg) marshal(out *ImBuffer) {

}

func (msg *CTestMsg) unmarshal(in *ImBuffer) {
}

func registerMsg(id int32, elem interface{}) {
	t := reflect.TypeOf(elem).Elem()
	msgs[id] = t
}

func init() {
	//registerMsg(123, (*CTestMsg)(nil))

	registerMsg(1, (*SChallenge)(nil))

}

func main3() {
	fmt.Println(" msg ", msgs)
}
