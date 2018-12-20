package common

type IMarshal interface {
	Marshal(octets *Octets)
	Unmarshal(octets *Octets)
	New() IMarshal
}

type Protocol interface {
	Marshal(octets *Octets)
	Unmarshal(octets *Octets)
	New() IMarshal
	GetMsgId() int32
}

type CfgObj interface {
	//Load(octets *Octets)
	//New() CfgObj
	GetTypeId() int32
}
