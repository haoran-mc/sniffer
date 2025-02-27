package mock

import (
	"fmt"
	"log"
	"net"

	"github.com/buger/goreplay/proto"
	"github.com/haoran-mc/sniffer/cache"
)

const BUF_SIZE int = 1024

func MockServerStart() {
	listener, err := net.Listen("tcp", "127.0.0.1:9523")
	if err != nil {
		log.Fatal("fail to listen: ", err.Error())
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Fail to accept connection:", err.Error())
			continue
		}

		// 监听到一个客户连接，处理
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, BUF_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Fail to read data:", err.Error())
		return
	}
	req := buf[:n]

	cacheId := string(proto.Header(req, []byte("X-SnifferId")))
	resp, ok := cache.GetResponse(cacheId)
	if ok {
		_, err = conn.Write([]byte(resp))
		if err != nil {
			fmt.Println("Fail to send response:", err.Error())
		}
	}
}
