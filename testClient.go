package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"sync"
)

type ClientConnection struct {
	conn   *net.TCPConn
	userId int64
	state  int
}

func main() {
	fmt.Println(" start client !!!!")

	var buf [512]byte
	//if len(os.Args) != 2 {
	//	fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
	//	os.Exit(1)
	//}
	//service := os.Args[1]
	//service := "10.95.134.20:31004"
	service := "127.0.0.1:31004"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	checkErr(err)
	rAddr := conn.RemoteAddr()
	//n, err := conn.Write([]byte("Hello server!"))
	//checkErr(err)

	client.conn = conn

	n, err := client.conn.Read(buf[0:])

	for ; n > 0; {
		checkErr(err)
		fmt.Printf(" read data n = %v, err = %v \n", n, err)

		rev := make([]byte, n)
		copy(rev, buf[0:n])

		fmt.Println(" rev data from server ", rAddr.String(), rev, len(rev), cap(rev))
		//fmt.Println("Reply from server ", rAddr.String(), string(buf[0:n]))

		revBuffer := ImBuffer{rev, 0, 0}
		msgSize := revBuffer.readInt()
		msgType := revBuffer.readInt()
		fmt.Println(" rev msg ", msgSize, msgType)

		message := msgIdMap[msgType].New()
		message.decode(&revBuffer)

		fmt.Println("===== rev msg type ", reflect.TypeOf(message))
		handler, _ := msgHandlers.Load(msgType)
		fmt.Println(" rev msg ", msgSize, msgType, handler, reflect.TypeOf(message))

		//go func() { handler.(func(client *ClientConnection, msg Message))(client, message) }()
		handler.(func(client *ClientConnection, msg Message))(client, message)

		//
		//switch msgType {
		//case 1:
		//	msg := SChallenge{}
		//	msg.decode(&revBuffer)
		//	fmt.Println("rece msg ", msg, msg.nonce)
		//
		//	auth := buildAuth("wangbo", msg.nonce)
		//	n, err = conn.Write(auth)
		//	fmt.Println(" send data to server ", rAddr.String(), auth, n, err)
		//	checkErr(err)
		//}
		//msg := GAnnouceServerInfo{}
		//msg.decode(&revBuffer)
		//fmt.Println("rece msg ", msg)

		fmt.Println(" after send ")
		n, err = client.conn.Read(buf[0:])
	}

	fmt.Printf(" finish n = %v, err = %v \n", n, err)
	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

type SChallenge struct {
	nonce *ImBuffer
}

func (msg *SChallenge) encode(buffer *ImBuffer) {
	buffer.writeBuffer(msg.nonce)
}

func (msg *SChallenge) decode(buffer *ImBuffer) {
	msg.nonce = buffer.readBuffer()

	fmt.Println(" rev nonce", msg.nonce)
}

func (msg *SChallenge) New() Message {
	return &SChallenge{}
}

func (msg *SChallenge) getMsgId() int32 {
	return 1
}

func buildAuth(name string, nonce *ImBuffer) []byte {
	msg := CAuth{name, nonce}
	buffer := ImBuffer{}
	buffer.writeInt(2) //msg id = 2
	msg.encode(&buffer)

	fmt.Println(" auth length ", buffer)

	msgLength := buffer.length()

	fmt.Println(" auth length ", msgLength)
	buffer2 := ImBuffer{}
	//buffer2.writeMsgSize(msgLength)
	buffer2.writeBytes(buffer.data)
	//buffer2.writeBuffer(&buffer)

	fmt.Println(" auth all length ", len(buffer2.data), buffer2.data)

	return buffer2.data
}

type GAnnouceServerInfo struct {
	typeId   int32
	serverId int32
	extra    string
}

func (msg *GAnnouceServerInfo) encode(buffer *ImBuffer) {
	buffer.writeInt(msg.typeId)
	buffer.writeInt(msg.serverId)
	buffer.writeString(msg.extra)
}

func (msg *GAnnouceServerInfo) decode(buffer *ImBuffer) {
	msg.typeId = buffer.readInt()
	msg.serverId = buffer.readInt()

	fmt.Println(" GAnnouceServerInfo ", msg.typeId, msg.serverId)
	msg.extra = buffer.readString()
}

func getMessage(msgId int32) Message {
	switch msgId {
	case 1:
		return &SChallenge{}
	case 2:
		return &CAuth{}
	default:
		panic(" no such msgId " + string(msgId))
	}
}

type MsgHandle func(client ClientConnection, msg interface{})

var (
	msgIdMap    = make(map[int32]Message, 100)
	msgHandlers = sync.Map{}
	client      = &ClientConnection{state: LINKING}

	LINKING  = 0
	LINKED   = 1
	LOGINING = 2
	LOGINED  = 3
)

type CAuth struct {
	username string
	nonce    *ImBuffer
}

func (msg *CAuth) encode(buffer *ImBuffer) {
	buffer.writeString(msg.username)
	buffer.writeBuffer(msg.nonce)
}

func (msg *CAuth) decode(buffer *ImBuffer) {
	msg.username = buffer.readString()
	msg.nonce = buffer.readBuffer()
}

func (msg *CAuth) New() Message {
	return &CAuth{}
}

func (msg *CAuth) getMsgId() int32 {
	return 2
}

func handleSChallenge(client *ClientConnection, msg Message) {
	challenge := msg.(*SChallenge)
	fmt.Println(" handleSChallenge ", msg)

	fmt.Println("rece msg ", msg, challenge.nonce)

	retMsg := CAuth{"wangbo", challenge.nonce}

	sendMessage(client, &retMsg)

	//auth := buildAuth("wangbo", challenge.nonce)
	//n, err := client.conn.Write(auth)
	//fmt.Println(" send data to server ", client.conn.RemoteAddr().String(), auth, n, err)
	//checkErr(err)
}

func handleCAuth(client ClientConnection, msg Message) {
	fmt.Println(" handleCAuth ", msg)
}

func handleSAuth(client *ClientConnection, msg Message) {
	auth := msg.(*SAuth)

	fmt.Println(" handleSAuth ", auth.username, auth.userid, auth.err)

	if auth.err == 0 {
		client.userId = auth.userid
		fmt.Println(" get user Id ", client.userId)
		retMsg := CBindServer{30001}
		sendMessage(client, &retMsg)
	}
}

func handleSBindServer(client *ClientConnection, msg Message) {
	bindServer := msg.(*SBindServer)
	fmt.Println(" bind server ret ", bindServer.err, bindServer.serverId)

	fmt.Println(" send login msg ", client.userId)

	if bindServer.err == 0 {
		client.state = LINKED
		sendMessage(client, &CLogin{client.userId})
	}
}

func handleKeepAlive(client *ClientConnection, msg Message) {
	alive := msg.(*KeepAlive)
	fmt.Println(" handleKeepAlive ", alive)

}

func sendMessage(client *ClientConnection, msg Message) {

	var m1 = msg
	if client.state >= LINKED {
		gsMsg := CForward{}
		gsMsg.msgId = msg.getMsgId()
		buffer := ImBuffer{}
		msg.encode(&buffer)
		gsMsg.data = &buffer

		m1 = &gsMsg
	}

	//auth := buildAuth("wangbo", challenge.nonce)

	//msg := CAuth{name, nonce}
	buffer := ImBuffer{}
	buffer.writeInt(m1.getMsgId())
	m1.encode(&buffer)

	//fmt.Println(" auth length ", buffer)

	//msgLength := buffer.length()

	//fmt.Println(" auth length ", msgLength)
	buffer2 := ImBuffer{}
	//buffer2.writeMsgSize(msgLength)
	buffer2.writeBytes(buffer.data)
	//buffer2.writeBuffer(&buffer)

	//fmt.Println(" auth all length ", len(buffer2.data), buffer2.data)

	n, err := client.conn.Write(buffer2.data)
	fmt.Printf(" send data to server msg = %v %v, write bytes =  %v, err = %v \n", reflect.TypeOf(msg), msg, n, err)
	checkErr(err)
}

func init() {

	allMessages := []Message{&SChallenge{}, &CAuth{}, &SAuth{}, &CBindServer{}, &SBindServer{}, &CLogin{}, &KeepAlive{}, &CForward{}}

	for _, m := range allMessages {
		msgIdMap[m.getMsgId()] = m
	}

	msgHandlers.Store(int32(0), handleKeepAlive)
	msgHandlers.Store(int32(1), handleSChallenge)
	msgHandlers.Store(int32(2), handleCAuth)
	msgHandlers.Store(int32(3), handleSAuth)
	msgHandlers.Store(int32(7), handleSBindServer)

	value, ok := msgHandlers.Load(1)
	fmt.Println(" init msg handler map ", msgHandlers, value, ok)
	fmt.Println(" init msg id map ", msgIdMap)

}

type Message interface {
	encode(buffer *ImBuffer)
	decode(buffer *ImBuffer)
	New() Message
	getMsgId() int32
}

type SAuth struct {
	err      int32
	username string
	userid   int64
}

func (msg *SAuth) encode(buffer *ImBuffer) {
	buffer.writeInt(msg.err)
	buffer.writeString(msg.username)
	buffer.writeLong(msg.userid)
}

func (msg *SAuth) decode(buffer *ImBuffer) {
	msg.err = buffer.readInt()
	msg.username = buffer.readString()
	msg.userid = buffer.readLong()
}

func (msg *SAuth) New() Message {
	return &SAuth{}
}

func (msg *SAuth) getMsgId() int32 {
	return 3
}

type CBindServer struct {
	serverId int32
}

func (msg *CBindServer) getMsgId() int32 {
	return 6
}

func (msg *CBindServer) encode(buffer *ImBuffer) {
	buffer.writeInt(msg.serverId)
}

func (msg *CBindServer) decode(buffer *ImBuffer) {
	msg.serverId = buffer.readInt()
}

func (msg *CBindServer) New() Message {
	return &CBindServer{}
}

type SBindServer struct {
	err      int32
	serverId int32
}

func (msg *SBindServer) encode(buffer *ImBuffer) {
	buffer.writeInt(msg.err)
	buffer.writeInt(msg.serverId)
}

func (msg *SBindServer) decode(buffer *ImBuffer) {
	msg.err = buffer.readInt()
	msg.serverId = buffer.readInt()
}

func (msg *SBindServer) New() Message {
	return &SBindServer{}
}

func (*SBindServer) getMsgId() int32 {
	return 7
}

type CLogin struct {
	userId int64
}

func (msg *CLogin) encode(buffer *ImBuffer) {
	buffer.writeLong(msg.userId)
}

func (msg *CLogin) decode(buffer *ImBuffer) {
	msg.userId = buffer.readLong()
}

func (msg *CLogin) New() Message {
	return &CLogin{}
}

func (msg *CLogin) getMsgId() int32 {
	return 60768
}

type KeepAlive struct {
}

func (*KeepAlive) encode(buffer *ImBuffer) {
}

func (*KeepAlive) decode(buffer *ImBuffer) {
}

func (*KeepAlive) New() Message {
	return &KeepAlive{}
}

func (*KeepAlive) getMsgId() int32 {
	return 0
}

//// gen code
type CForward struct {
	msgId int32
	data  *ImBuffer
}

func (msg *CForward) encode(buffer *ImBuffer) {
	buffer.writeInt(msg.msgId)
	buffer.writeBuffer(msg.data)
}

func (msg *CForward) decode(buffer *ImBuffer) {
	msg.msgId = buffer.readInt()
	msg.data = buffer.readBuffer()
}

func (msg *CForward) New() Message {
	return &CForward{0, &ImBuffer{}}
}

func (msg *CForward) getMsgId() int32 {
	return 4
}
