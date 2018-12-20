package main

import "awesomeProject/common"

type KeepAlive struct {
}

func (*KeepAlive) Marshal(octets *common.Octets) {
}

func (*KeepAlive) Unmarshal(octets *common.Octets) {
}

func (*KeepAlive) New() common.IMarshal {
	return &KeepAlive{}
}

func (*KeepAlive) GetMsgId() int32 {
	return 0
}
