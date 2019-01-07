package main

import (
	"net"
	"sync"
	"bufio"
	"os"
	"go-min-chat/Utils"
	"strings"
	"github.com/golang/protobuf/proto"
	"github.com/weiwenwang/ApiBench/protobuf/proto"
	"github.com/weiwenwang/ApiBench/client/cli"
	"github.com/weiwenwang/ApiBench/const"
	"github.com/weiwenwang/DiyProtocol"
	"encoding/binary"
	"bytes"
	"fmt"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8089")
	//go
	//buf := make([]byte, 50)
	//conn.Read(buf)
	//fmt.Println(buf)
	cli := cli.GetCli()
	cli.Pattern = "DEBUG"
	var wg sync.WaitGroup
	ch := make(chan []byte)
	wg.Add(3)
	go readFromStdio(ch)
	// 接收消息
	go readFromConn(conn)
	// 发送消息
	go Send(conn, ch)
	fmt.Printf("127.0.0.1:8089>")
	// 发送登录的消息
	//go heartBeat(conn)
	// todo 这里是一直等，要去解决
	wg.Wait()

}

func readFromStdio(ch chan []byte) {
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		data = []byte(strings.Trim(string(data), " "))
		data_str := string(data)
		if (strings.EqualFold(data_str, "")) { // 直接按的回车, 不做处理
			continue
		}
		p1 := &protobuf.Content{}
		p1.Id = 1;
		p1.Command = data_str
		d, _ := proto.Marshal(p1)
		ch <- d
	}
}

func readFromConn(conn net.Conn) {
	msg := make(chan string, 10)
	for {
		select {
		case msg1 := <-msg: // 从msg chan里面取数据
			doMsg(msg1)
		default: // 读到数据放到msg chan里面
			buffer := DiyProtocol.NewBuffer(conn, "BEGIN", 1024)
			buffer.Read(msg)
		}
	}
}

func doMsg(buf string) {
	backContent := &protobuf.BackContent{}
	proto.Unmarshal([]byte(buf), backContent)
	switch backContent.Id {
	case _const.LIST:
		Utils.EchoLine("LIST", 2)
		break
	case _const.UNDOLIST:
		Utils.EchoLine("UNDOLIST", 1)
		break
	case _const.DONELIST:
		Utils.EchoLine("DONELIST", 1)
		break
	case _const.INFO:
		Utils.EchoLine("INFO", 1)
		break
	}
}

func Send(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		send(conn, content)
	}
}

func send(conn net.Conn, ch []byte) {
	fmt.Println("发送了")
	headSize := len(ch)
	var headBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(headBytes, uint16(headSize))
	var buffer_client bytes.Buffer
	buffer_client.Write([]byte("BEGIN"))
	buffer_client.Write(headBytes)
	buffer_client.Write(ch)
	b3 := buffer_client.Bytes() //得到了b1+b2的结果
	conn.Write(b3)
}
