package main

import (
	"fmt"
	"net"
	"go-min-chat/Utils"
	"github.com/weiwenwang/DiyProtocol"
	"github.com/weiwenwang/ApiBench/protobuf/proto"
	"github.com/weiwenwang/ApiBench/const"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"bytes"
	"strings"
	"github.com/weiwenwang/ApiBench/common"
	"github.com/weiwenwang/ApiBench/server/ser"
)

func main() {
	listen, _ := net.Listen("tcp", "127.0.0.1:8089")
	for {
		newConn, err := listen.Accept()
		// 连接上了，就把这个连接赋予一个未登录的用户
		fmt.Println(newConn.RemoteAddr())
		Utils.CheckError(err)
		msg := make(chan []byte, 100)
		// 收到消息
		go recvConnMsg(newConn, msg)
		// 发送消息
		go sendConnMsg(newConn, msg)
	}
}

// 服务端发送消息
func sendConnMsg(conn net.Conn, ch chan []byte) {
	for {
		content, _ := <-ch
		send(conn, content)
	}
}

// 服务端接收消息
func recvConnMsg(conn net.Conn, msg_chan chan []byte) {
	msg1 := make(chan string, 20)                         // 这里设置消息channel可以容纳10个消息
	buffer1 := DiyProtocol.NewBuffer(conn, "BEGIN", 1024) // 缓存区设置1024字节， 如果单个消息大于这个值就不能接受了
	for {
		select {
		case msg2 := <-msg1: // 从msg chan里面取数据
			doMsg(msg2, msg_chan)
			// 读到数据放到msg chan里面, 如果msg里没有消息了，那么会阻塞在ThomasRead这个函数里等消息
			// 当读到数据了，并分析这次读到的数据里有几个消息, 不足一个消息就继续读，大于一个消息了，处理掉这个消息
		default:
			buffer1.Read(msg1)
		}
	}
}

func doMsg(msg string, msg_chan chan []byte) {
	Content := &protobuf.Content{}
	proto.Unmarshal([]byte(msg), Content)
	fmt.Println(Content.Command)
	//addrequest url http://www.baidu.com -c 200 -s 1000
	param := strings.Split(Content.Command, " ")
	GetMsgType(param)
	if Content.Command == "List" {
		p_back := &protobuf.BackContent{}
		p_back.Id = _const.LIST
		p_back.Msg = "list"
		data, _ := proto.Marshal(p_back)
		msg_chan <- data
	}
}

func GetMsgType(param []string) {
	ser1 := ser.GetSer()
	if param[0] == "addrequest" {
		re := common.Request{}
		re.C = 200
		re.Ps = 10
		// 添加到列表里面去
		ser1.List = append(ser1.List, re)
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
