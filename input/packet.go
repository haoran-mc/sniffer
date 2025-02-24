package input

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type packet struct {
	data []byte
}

var pktChan chan *packet

func init() {
	pktChan = make(chan *packet, 1000) // maximum number of packets in the receive queue
	go func() {
		for {
			pkt := <-pktChan
			tcpIpPkt, err := pkt.extractTcpPacket()
			if err != nil {
				fmt.Println("extract packet failed:", err.Error())
				continue
			}

			// tcpReassemble TODO TCP Segmentation
			tcpMessage := TcpMessage{
				packet: tcpIpPkt,
			}
			tcpMessage.processPacket(tcpIpPkt)
			if tcpMessage.Protocol == Http {
				// send request
				// receive response
			}
		}
	}()
}

func (pkt *packet) extractTcpPacket() (tcpIpPkt *tcpIpPacket, err error) {
	if len(pkt.data) <= 14 { // ipv4 or ipv6
		return nil, errors.New("eth data length is invalid")
	}

	tcpIpPkt = new(tcpIpPacket)
	var ipPkt = pkt.data[14:] // ethernet header length is 14
	var tcpSeg = []byte{}

	// IP Layer
	switch ipPkt[0] >> 4 {
	case 4:
		ihl := int((ipPkt[0] & 0x0F) * 4)
		if ihl < 20 || len(ipPkt) < ihl {
			return nil, errors.New("ipv4 header length is invalid")
		}
		tl := int(ipPkt[2])<<8 + int(ipPkt[3])
		if len(ipPkt) < tl {
			return nil, errors.New("ipv4 packet length is invalid")
		}
		if ipPkt[9] != 0x06 { // TCP: 06
			return nil, errors.New("not tcp")
		}
		tcpIpPkt.IsIPv6 = false
		tcpIpPkt.SrcIP = ipPkt[12:16]
		tcpIpPkt.DstIP = ipPkt[16:20]
		tcpSeg = ipPkt[ihl:]
	case 6:
		if len(ipPkt) < 40 {
			return nil, errors.New("ipv6 header length is invalid")
		}
		if ipPkt[6] != 0x06 { // TCP: 06
			return nil, errors.New("not tcp")
		}
		tcpIpPkt.IsIPv6 = true
		tcpIpPkt.SrcIP = ipPkt[8:24]
		tcpIpPkt.DstIP = ipPkt[24:40]
		tcpSeg = ipPkt[40:]
	default:
		return nil, errors.New("unsupported IP version")
	}

	// TCP Layer
	if len(tcpSeg) < 20 {
		return nil, errors.New("tcp header length is invalid")
	}
	dataOffset := int((tcpSeg[12] >> 4) * 4)
	if dataOffset < 20 || len(tcpSeg) < dataOffset {
		return nil, errors.New("tcp header length is invalid")
	}
	srcPort := int(tcpSeg[0])<<8 + int(tcpSeg[1])
	dstPort := int(tcpSeg[2])<<8 + int(tcpSeg[3])
	tcpIpPkt.SrcPort = uint16(srcPort)
	tcpIpPkt.DstPort = uint16(dstPort)
	tcpIpPkt.Seq = binary.BigEndian.Uint32(tcpSeg[4:8])
	tcpIpPkt.Ack = binary.BigEndian.Uint32(tcpSeg[8:12])
	tcpIpPkt.URG = (tcpSeg[13] & 0x20) != 0 // 00100000
	tcpIpPkt.ACK = (tcpSeg[13] & 0x10) != 0
	tcpIpPkt.PSH = (tcpSeg[13] & 0x08) != 0
	tcpIpPkt.RST = (tcpSeg[13] & 0x04) != 0
	tcpIpPkt.SYN = (tcpSeg[13] & 0x02) != 0
	tcpIpPkt.FIN = (tcpSeg[13] & 0x01) != 0
	tcpIpPkt.Window = binary.BigEndian.Uint16(tcpSeg[14:16])
	tcpIpPkt.CheckSum = tcpSeg[16:18]
	tcpIpPkt.Payload = tcpSeg[dataOffset:]
	return
}
