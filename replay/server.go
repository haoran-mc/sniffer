package replay

import (
	"fmt"
	"log"
	"net"

	"github.com/buger/goreplay/proto"
	"github.com/haoran-mc/sniffer/cache"
)

const bufSize = 1024

func StartResponseServer() {
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, bufSize)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Fail to read data:", err.Error())
		return
	}
	req := buf[:n]

	cacheID := string(proto.Header(req, []byte("X-SnifferId")))
	resp, ok := cache.GetResponse(cacheID)
	if ok {
		_, err = conn.Write([]byte(resp))
		if err != nil {
			fmt.Println("Fail to send response:", err.Error())
		}
	}
}
