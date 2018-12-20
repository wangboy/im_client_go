package main

import (
	"awesomeProject/common"
	"awesomeProject/gen/cfg/cfg/skill"
	"awesomeProject/gen/msg/msg/gs/login"
	. "awesomeProject/gen/msg/msg/link"
	"fmt"
	"net"
	"os"
	"path"
	"reflect"
	"sync"
)

type ImClient struct {
	conn   *net.TCPConn
	userId int64
	state  int
}

func main() {

	var array = [3]int{0, 1, 2}
	var array2 = &array

	array2[2] = 5
	fmt.Println(array, *array2)

	testLoadFile()

	TestCfg()

	startClient()
}

func testLoadFile() {
	filePath := path.Join("gen/config/", "buffcfg.data")
	common.BinRead(filePath)
}

func TestCfg() {
	cfg := skill.SkillCfg{}
	fmt.Println(cfg)
}

func startClient() {

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

	for n > 0 {
		checkErr(err)
		fmt.Printf(" read Data n = %v, err = %v \n", n, err)

		rev := make([]byte, n)
		copy(rev, buf[0:n])

		fmt.Println("===== rev Data from server ", rAddr.String(), rev, len(rev), cap(rev))
		//fmt.Println("Reply from server ", rAddr.String(), string(buf[0:n]))

		revBuffer := common.Octets{rev, 0, 0}
		msgSize := revBuffer.ReadInt()
		msgType := revBuffer.ReadInt()
		fmt.Println(" rev msg ", msgSize, msgType)

		message := msgIdMap[msgType].New()
		message.Unmarshal(&revBuffer)

		fmt.Println("rev msg type ", reflect.TypeOf(message))
		handler, _ := msgHandlers.Load(msgType)
		fmt.Println(" rev msg ", msgSize, msgType, handler, reflect.TypeOf(message))

		go func() { handler.(func(client *ImClient, msg common.IMarshal))(client, message) }()
		//handler.(func(client *ImClient, msg common.IMarshal))(client, message)

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

type MsgHandle func(client ImClient, msg interface{})

var (
	msgIdMap    = make(map[int32]common.IMarshal, 100)
	msgHandlers = sync.Map{}

	client = &ImClient{state: LINKING}

	LINKING  = 0
	LINKED   = 1
	LOGINING = 2
	LOGINED  = 3
)

func sendMessage(client *ImClient, msg common.Protocol) {

	var m1 = msg
	if client.state >= LINKED {
		gsMsg := CForward{}
		gsMsg.Type_ = msg.GetMsgId()
		buffer := common.Octets{}
		msg.Marshal(&buffer)
		gsMsg.Data_ = buffer

		m1 = &gsMsg
	}

	buffer := common.Octets{}
	buffer.WriteInt(m1.GetMsgId())
	m1.Marshal(&buffer)

	buffer2 := common.Octets{}
	buffer2.WriteBytes(buffer.Data)

	n, err := client.conn.Write(buffer2.Data)
	fmt.Printf(" send Data to server msg = %v %v, write bytes =  %v, err = %v \n", reflect.TypeOf(msg), msg, n, err)
	checkErr(err)
}

func init() {

	allMessages := []common.Protocol{&SChallenge{}, &CAuth{}, &SAuth{}, &CBindServer{}, &SBindServer{}, &login.CLogin{}, &KeepAlive{}, &CForward{}}

	for _, m := range allMessages {
		msgIdMap[m.GetMsgId()] = m
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

func handleSChallenge(client *ImClient, msg common.IMarshal) {
	challenge := msg.(*SChallenge)
	fmt.Println(" handleSChallenge ", msg)

	fmt.Println("rece msg ", msg, challenge.Nonce_)

	retMsg := CAuth{"wangbo", challenge.Nonce_}

	sendMessage(client, &retMsg)

}

func handleCAuth(client ImClient, msg common.IMarshal) {
	fmt.Println(" handleCAuth ", msg)
}

func handleSAuth(client *ImClient, msg common.IMarshal) {
	auth := msg.(*SAuth)

	fmt.Println(" handleSAuth ", auth.Username_, auth.Userid_, auth.Err_)

	if auth.Err_ == 0 {
		client.userId = auth.Userid_
		fmt.Println(" get user Id ", client.userId)
		retMsg := CBindServer{30001}
		sendMessage(client, &retMsg)
	}
}

func handleSBindServer(client *ImClient, msg common.IMarshal) {
	bindServer := msg.(*SBindServer)
	fmt.Println(" bind server ret ", bindServer.Err_, bindServer.Serverid_)

	fmt.Println(" send login msg ", client.userId)

	if bindServer.Err_ == 0 {
		client.state = LINKED
		sendMessage(client, &login.CLogin{client.userId})
	}
}

func handleKeepAlive(client *ImClient, msg common.IMarshal) {
	alive := msg.(*KeepAlive)
	fmt.Println(" handleKeepAlive ", alive)

}
