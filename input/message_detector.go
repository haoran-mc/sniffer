package input

type messageDetector struct{}

func (messageDetector) detect(flow *tcpFlow, pkt *tcpIpPacket) *TcpMessage {
	// 检测入口已经切到流上下文，后续可以在这里平滑接入真正的 TCP 重组。
	return detectTCPMessage(flow, pkt)
}

func detectTCPMessage(flow *tcpFlow, pkt *tcpIpPacket) *TcpMessage {
	if len(pkt.Payload) == 0 {
		return nil
	}

	candidatePayload := flow.candidatePayload(pkt)
	protocol, direction, ok := detectApplicationMessage(candidatePayload)
	if !ok {
		flow.remember(pkt)
		return nil
	}

	flow.reset(pkt)

	return &TcpMessage{
		packet:    pkt,
		payload:   candidatePayload,
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
