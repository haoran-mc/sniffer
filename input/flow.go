package input

import "fmt"

const maxFlowPreviewBytes = 8 * 1024

type packetFlowTracker struct {
	flows map[string]*tcpFlow
}

func newPacketFlowTracker() *packetFlowTracker {
	return &packetFlowTracker{
		flows: make(map[string]*tcpFlow),
	}
}

func (t *packetFlowTracker) get(pkt *tcpIpPacket) *tcpFlow {
	key := buildFlowKey(pkt)
	flow, ok := t.flows[key]
	if ok {
		return flow
	}

	flow = newTCPFlow(key, pkt)
	t.flows[key] = flow
	return flow
}

func (t *packetFlowTracker) release(pkt *tcpIpPacket) {
	if !(pkt.FIN || pkt.RST) {
		return
	}

	delete(t.flows, buildFlowKey(pkt))
}

type tcpFlow struct {
	key         string
	client      tcpEndpoint
	server      tcpEndpoint
	requestBuf  tcpStreamBuffer
	responseBuf tcpStreamBuffer
}

func newTCPFlow(key string, pkt *tcpIpPacket) *tcpFlow {
	return &tcpFlow{
		key:    key,
		client: newTCPEndpoint(pkt.SrcIP.String(), pkt.SrcPort),
		server: newTCPEndpoint(pkt.DstIP.String(), pkt.DstPort),
	}
}

func (f *tcpFlow) candidatePayload(pkt *tcpIpPacket) []byte {
	return f.bufferFor(pkt).candidate(pkt.Payload)
}

func (f *tcpFlow) remember(pkt *tcpIpPacket) {
	f.bufferFor(pkt).append(pkt.Payload)
}

func (f *tcpFlow) reset(pkt *tcpIpPacket) {
	f.bufferFor(pkt).reset()
}

func (f *tcpFlow) bufferFor(pkt *tcpIpPacket) *tcpStreamBuffer {
	if f.isClientToServer(pkt) {
		return &f.requestBuf
	}

	return &f.responseBuf
}

func (f *tcpFlow) isClientToServer(pkt *tcpIpPacket) bool {
	return pkt.SrcPort == f.client.Port && pkt.DstPort == f.server.Port &&
		pkt.SrcIP.String() == f.client.IP && pkt.DstIP.String() == f.server.IP
}

type tcpStreamBuffer struct {
	prefix []byte
}

func (b *tcpStreamBuffer) candidate(payload []byte) []byte {
	if len(b.prefix) == 0 {
		return payload
	}

	candidate := make([]byte, 0, len(b.prefix)+len(payload))
	candidate = append(candidate, b.prefix...)
	candidate = append(candidate, payload...)
	return candidate
}

func (b *tcpStreamBuffer) append(payload []byte) {
	if len(payload) == 0 {
		return
	}

	b.prefix = append(b.prefix, payload...)
	if len(b.prefix) > maxFlowPreviewBytes {
		b.prefix = b.prefix[len(b.prefix)-maxFlowPreviewBytes:]
	}
}

func (b *tcpStreamBuffer) reset() {
	b.prefix = nil
}

type tcpEndpoint struct {
	IP   string
	Port uint16
}

func newTCPEndpoint(ip string, port uint16) tcpEndpoint {
	return tcpEndpoint{IP: ip, Port: port}
}

func buildFlowKey(pkt *tcpIpPacket) string {
	left := fmt.Sprintf("%s:%d", pkt.SrcIP.String(), pkt.SrcPort)
	right := fmt.Sprintf("%s:%d", pkt.DstIP.String(), pkt.DstPort)
	if left < right {
		return left + "-" + right
	}

	return right + "-" + left
}
