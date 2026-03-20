package input

import (
	"encoding/binary"
	"encoding/hex"
	"net"
)

func buildMessageUUID(pkt *tcpIpPacket, direction Dir) []byte {
	var streamID uint64

	// 请求和响应使用统一的方向基准，确保能生成相同的关联 ID。
	if direction == DirIncoming {
		streamID = uint64(pkt.SrcPort)<<48 | uint64(pkt.DstPort)<<32 |
			uint64(ipToInt(pkt.SrcIP))
	} else {
		streamID = uint64(pkt.DstPort)<<48 | uint64(pkt.SrcPort)<<32 |
			uint64(ipToInt(pkt.DstIP))
	}

	id := make([]byte, 12)
	binary.BigEndian.PutUint64(id, streamID)

	if direction == DirIncoming {
		binary.BigEndian.PutUint32(id[8:], pkt.Ack)
	} else {
		binary.BigEndian.PutUint32(id[8:], pkt.Seq)
	}

	uuidHex := make([]byte, 24)
	hex.Encode(uuidHex[:], id[:])

	return uuidHex
}

func ipToInt(ip net.IP) uint32 {
	if len(ip) == 0 {
		return 0
	}

	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}

	return binary.BigEndian.Uint32(ip)
}
