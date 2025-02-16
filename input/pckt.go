package input

import (
	"fmt"

	"github.com/google/gopacket"
)

type pckt struct {
	data        []byte
	captureInfo *gopacket.CaptureInfo
}

var pcktChan chan *pckt

func init() {
	pcktChan = make(chan *pckt, 1000) // maximum number of packets in the receive queue
	go func() {
		for {
			pkt := <-pcktChan
			fmt.Println(pkt)
			// 1. 判断是否为 tcp 包
			// 2. 解析 tcp 包
		}
	}()
}
