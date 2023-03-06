package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"stress/src"
	"stress/src/protoA"
	"time"

	"google.golang.org/protobuf/proto"
)

var (
	tm = map[uint32]proto.Message{
		1001: &protoA.LoginReq{},
		1002: &protoA.LoginResp{},
		1004: &protoA.SayHelloResp{},
		1003: &protoA.SayHelloReq{},
	}
)

func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	defer func() {
		recover()
	}()
	codec := src.TestCodec{}
	t := 100
	for {

		id, data, _ := codec.Decode(conn)

		if resp, ok := tm[uint32(id)]; ok {
			proto.Unmarshal(data, resp)
			fmt.Println("收到Client端发来的数据：", resp)
			// time.Sleep(time.Second * 200)
			time.Sleep(time.Duration(t) * time.Millisecond)
			// t += 100
			if id == 1001 {
				data := &protoA.LoginResp{Greet: "发送的数据"}
				dd, _ := proto.Marshal(data)
				cc, _ := codec.Encode(1002, dd)
				conn.Write(cc)
			} else {
				data := &protoA.SayHelloResp{Name: "aaaa"}
				dd, _ := proto.Marshal(data)
				cc, _ := codec.Encode(1004, dd)
				conn.Write(cc)
			}

		} else {
			fmt.Println("收到不匹配：", id)
			return
		}

		// conn.Write([]byte(recvStr)) // 发送数据
	}
}

func Server() {
	listen, err := net.Listen("tcp", src.ServerAddr)
	if err != nil {
		fmt.Println("Listen() failed, err: ", err)
		return
	}
	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		fmt.Println("jbk")
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		go process(conn) // 启动一个goroutine来处理客户端的连接请求
	}

}

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	Server()
}
