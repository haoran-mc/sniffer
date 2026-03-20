package input

type packetProcessor struct {
	detector   messageDetector
	dispatcher messageDispatcher
}

func newPacketProcessor() *packetProcessor {
	// 这里先保留一个最小处理流水线：识别后直接分发。
	return &packetProcessor{
		detector:   messageDetector{},
		dispatcher: messageDispatcher{},
	}
}

func (p *packetProcessor) process(pkt *packet) {
	tcpIPPacket, err := pkt.extractTcpPacket()
	if err != nil {
		return
	}

	tcpMessage := p.detector.detect(tcpIPPacket)
	if tcpMessage == nil {
		return
	}

	p.dispatcher.dispatch(tcpMessage)
}
