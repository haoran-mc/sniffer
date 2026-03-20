package input

type packetProcessor struct {
	detector   messageDetector
	dispatcher messageDispatcher
}

func newPacketProcessor() *packetProcessor {
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

type messageDetector struct{}

func (messageDetector) detect(pkt *tcpIpPacket) *TcpMessage {
	// tcpReassemble TODO TCP Segmentation
	tcpMessage := &TcpMessage{
		packet: pkt,
	}
	tcpMessage.processPacket(pkt)
	if tcpMessage.Protocol == Unknown {
		return nil
	}

	return tcpMessage
}
