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

