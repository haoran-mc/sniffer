package cache

var (
	requestCache  map[string][]byte
	responseCache map[string][]byte
)

func init() {
	requestCache = make(map[string][]byte)
	responseCache = make(map[string][]byte)
}

func GetRequest(cacheId string) ([]byte, bool) {
	req, get := requestCache[cacheId]
	return req, get
}

func SetRequest(cacheId string, request []byte) {
	requestCache[cacheId] = request
}

func DelRequest(cacheId string) {
	delete(requestCache, cacheId)
}

func GetResponse(cacheId string) ([]byte, bool) {
	resp, get := responseCache[cacheId]
	return resp, get
}

func SetResponse(cacheId string, response []byte) {
	responseCache[cacheId] = response
}

func DelResponse(cacheId string) {
	delete(responseCache, cacheId)
}
