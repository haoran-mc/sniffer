package input

import (
	"net"
)

type tcpIpPacket struct {
	// IP
	IsIPv6       bool
	SrcIP, DstIP net.IP

	// TCP
	SrcPort, DstPort             uint16
	Ack, Seq                     uint32
	URG, ACK, PSH, RST, SYN, FIN bool
	Window                       uint16
	CheckSum                     []byte
	Payload                      []byte
}

type msgType = uint8

const (
	Unknow msgType = iota
	HttpRequest
	HttpResponse
)

// detect protocol and direction
func detectMsgType(pkt *tcpIpPacket) msgType {
	if len(pkt.Payload) > 0 {
		switch {
		case HasRequestTitle(pkt.Payload):
			return HttpRequest
		case HasResponseTitle(pkt.Payload):
			return HttpResponse
		}
	}
	return Unknow
}
