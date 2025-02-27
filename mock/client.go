package mock

import (
	"bytes"
	"log"
	"net"
	"strconv"

	"github.com/buger/goreplay/proto"
)

func SendRequest(request []byte) {
	// 修正 Content-Length
	if contentLenStr := proto.Header(request, []byte("Content-Length")); len(contentLenStr) > 0 {
		contentLen, _ := strconv.Atoi(string(contentLenStr))
		bodySize := len(request) - bytes.Index(request, []byte("\r\n\r\n")) - 4
		if bodySize != contentLen {
			request = proto.SetHeader(request, []byte("Content-Length"), []byte(strconv.Itoa(bodySize)))
		}
	}
	request = proto.DeleteHeader(request, []byte("Connection"))

	conn, err := net.Dial("tcp", "127.0.0.1:9522")
	if err != nil {
		log.Fatal("fail to dial: ", err.Error())
	}
	conn.Write(request)
	conn.Close()
}
