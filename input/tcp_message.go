package input

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
	packet  *tcpIpPacket
	payload []byte

	Protocol  ApplicationLayerType
	Direction Dir
	uuid      []byte
}

func (m *TcpMessage) Payload() []byte {
	if len(m.payload) > 0 {
		return m.payload
	}

	if m.packet == nil {
		return nil
	}

	return m.packet.Payload
}
