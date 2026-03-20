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
	packet *tcpIpPacket

	Protocol  ApplicationLayerType
	Direction Dir
	uuid      []byte
}
