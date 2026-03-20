package input

type packetProcessor struct {
	flows      *packetFlowTracker
	detector   messageDetector
	dispatcher messageDispatcher
}

func newPacketProcessor() *packetProcessor {
	// 这里先保留一个最小处理流水线：识别后直接分发。
	return &packetProcessor{
		flows:      newPacketFlowTracker(),
		detector:   messageDetector{},
		dispatcher: messageDispatcher{},
	}
}

func (p *packetProcessor) process(pkt *packet) {
	tcpIPPacket, err := pkt.extractTcpPacket()
	if err != nil {
		return
	}
	defer p.flows.release(tcpIPPacket)

	flow := p.flows.get(tcpIPPacket)
	tcpMessage := p.detector.detect(flow, tcpIPPacket)
	if tcpMessage == nil {
		return
	}

	p.dispatcher.dispatch(tcpMessage)
}
