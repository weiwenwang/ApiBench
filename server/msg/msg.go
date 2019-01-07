package msg

import (
	"net"
	"github.com/weiwenwang/ApiBench_Server/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"fmt"
)

func DoAllMsg(conn net.Conn, string2 string) {
	p1 := &protobuf.Content{}
	proto.Unmarshal([]byte(string2), p1)
	fmt.Println(p1)
}
