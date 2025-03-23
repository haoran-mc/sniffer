package cache

import "sync"

var (
	requestCache  sync.Map
	responseCache sync.Map
)

func GetRequest(cacheId string) ([]byte, bool) {
	req, get := requestCache.Load(cacheId)
	if !get {
		return nil, false
	}
	return req.([]byte), true
}

func SetRequest(cacheId string, request []byte) {
	requestCache.Store(cacheId, request)
}

func DelRequest(cacheId string) {
	requestCache.Delete(cacheId)
}

func GetResponse(cacheId string) ([]byte, bool) {
	resp, get := responseCache.Load(cacheId)
	if !get {
		return nil, false
	}
	return resp.([]byte), true
}

func SetResponse(cacheId string, response []byte) {
	responseCache.Store(cacheId, response)
}

func DelResponse(cacheId string) {
	responseCache.Delete(cacheId)
}
