package input

import (
	"github.com/buger/goreplay/proto"
	"github.com/haoran-mc/sniffer/cache"
	"github.com/haoran-mc/sniffer/mock"
)

type messageDispatcher struct{}

func (messageDispatcher) dispatch(message *TcpMessage) {
	cacheID := string(message.uuid)

	// 按请求方向和响应方向走不同的缓存与回放路径。
	switch message.Direction {
	case DirIncoming:
		dispatchRequest(cacheID, message.packet.Payload)
	case DirOutcoming:
		dispatchResponse(cacheID, message.packet.Payload)
	}
}

func dispatchRequest(cacheID string, payload []byte) {
	// 给回放链路补上内部关联 ID，便于在 mock 服务里取回响应。
	request := proto.AddHeader(payload, []byte("X-SnifferId"), []byte(cacheID))

	_, hasResponse := cache.GetResponse(cacheID)
	if hasResponse {
		mock.SendRequest(request)
		return
	}

	cache.SetRequest(cacheID, request)
}

func dispatchResponse(cacheID string, payload []byte) {
	request, hasRequest := cache.GetRequest(cacheID)
	if hasRequest {
		// 先回放请求，再缓存响应，维持现有处理顺序不变。
		mock.SendRequest(request)
		cache.SetResponse(cacheID, payload)
		cache.DelRequest(cacheID)
		return
	}

	cache.SetResponse(cacheID, payload)
}
