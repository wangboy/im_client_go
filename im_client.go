package main

import (
	"bufio"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/examples/chat/proto"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	"os"
	"reflect"
	"strings"

	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

var logger = golog.New("client")

func ReadConsole(callback func(string)) {

	for {

		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}

		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)

		callback(text)

	}

}

func Client() {

	addr := "10.95.134.20:30421"

	// 创建一个事件处理队列，整个客户端只有这一个队列处理事件，客户端属于单线程模型
	queue := cellnet.NewEventQueue()

	// 创建一个tcp的连接器，名称为client，连接地址为127.0.0.1:8801，将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	p := peer.NewGenericPeer("tcp.Connector", "im_client", addr, queue)

	// 设定封包收发处理的模式为tcp的ltv(Length-Type-Value), Length为封包大小，Type为消息ID，Value为消息内容
	// 并使用switch处理收到的消息
	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:
			logger.Debugln("client connected")
		case *cellnet.SessionClosed:
			logger.Debugln("client error")
		case *proto.ChatACK:
			logger.Infof("sid%d say: %s", msg.Id, msg.Content)
		}
	})

	// 开始发起到服务器的连接
	p.Start()

	// 事件队列开始循环
	queue.StartLoop()

	// 阻塞的从命令行获取聊天输入
	ReadConsole(func(str string) {

		p.(interface {
			Session() cellnet.Session
		}).Session().Send(&proto.ChatREQ{
			Content: str,
		})

	})
}

func init_client() {

	fmt.Println(" in in init codec ")

	imCodec := new(ImCodec)

	// 注册编码器
	codec.RegisterCodec(imCodec)

	fmt.Println(" init codec ", imCodec)

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("ImCodes"),
		Type:  reflect.TypeOf((*SChallenge)(nil)).Elem(),
		ID:    1,
	})
}
