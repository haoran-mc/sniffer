package input

import (
	"encoding/binary"
	"encoding/hex"
	"net"
)

type ApplicationLayerType uint8

const (
	Unknown ApplicationLayerType = iota
	Http
)

type Dir int

const (
	DirUnknown = iota
	DirIncoming
	DirOutcoming
)

type TcpMessage struct {
	packet *tcpIpPacket

	Protocol  ApplicationLayerType
	Direction Dir
	uuid      []byte
}

// detect protocol and direction
func (tcpMessage *TcpMessage) processPacket(pkt *tcpIpPacket) {
	if len(pkt.Payload) > 0 {
		switch {
		case HasRequestTitle(pkt.Payload):
			tcpMessage.Protocol = Http
			tcpMessage.Direction = DirIncoming
			tcpMessage.uuid = tcpMessage.UUID()
		case HasResponseTitle(pkt.Payload):
			tcpMessage.Protocol = Http
			tcpMessage.Direction = DirOutcoming
			tcpMessage.uuid = tcpMessage.UUID()
		}
	}
}

// UUID returns the UUID of a TCP request and its response.
func (m *TcpMessage) UUID() []byte {
	var streamID uint64
	pckt := m.packet

	// check if response or request have generated the ID before.
	if m.Direction == DirIncoming {
		streamID = uint64(pckt.SrcPort)<<48 | uint64(pckt.DstPort)<<32 |
			uint64(ip2int(pckt.SrcIP))
	} else {
		streamID = uint64(pckt.DstPort)<<48 | uint64(pckt.SrcPort)<<32 |
			uint64(ip2int(pckt.DstIP))
	}

	id := make([]byte, 12)
	binary.BigEndian.PutUint64(id, streamID)

	if m.Direction == DirIncoming {
		binary.BigEndian.PutUint32(id[8:], pckt.Ack)
	} else {
		binary.BigEndian.PutUint32(id[8:], pckt.Seq)
	}

	uuidHex := make([]byte, 24)
	hex.Encode(uuidHex[:], id[:])

	return uuidHex
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 0 {
		return 0
	}

	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}
