package input

type messageDetector struct{}

func (messageDetector) detect(pkt *tcpIpPacket) *TcpMessage {
	// 当前仍然是基于单个 TCP payload 做协议识别，后续再引入重组能力。
	return detectTCPMessage(pkt)
}

func detectTCPMessage(pkt *tcpIpPacket) *TcpMessage {
	if len(pkt.Payload) == 0 {
		return nil
	}

	protocol, direction, ok := detectApplicationMessage(pkt.Payload)
	if !ok {
		return nil
	}

	return &TcpMessage{
		packet:    pkt,
		Protocol:  protocol,
		Direction: direction,
		uuid:      buildMessageUUID(pkt, direction),
	}
}

func detectApplicationMessage(payload []byte) (ApplicationLayerType, Dir, bool) {
	switch {
	case HasRequestTitle(payload):
		return Http, DirIncoming, true
	case HasResponseTitle(payload):
		return Http, DirOutcoming, true
	default:
		return Unknown, DirUnknown, false
	}
}
