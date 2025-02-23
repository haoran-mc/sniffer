package input

import (
	"bytes"
	"net/http"
	"strings"
	"unsafe"
)

const (
	MinRequestCount  = 16 // GET / HTTP/1.1\r\n
	MinResponseCount = 14 // HTTP/1.1 200\r\n
	VersionLen       = 8  // HTTP/1.1
)

var CRLF = []byte("\r\n")

// Method returns HTTP method
func Method(payload []byte) []byte {
	end := bytes.IndexByte(payload, ' ')
	if end == -1 {
		return nil
	}

	return payload[:end]
}

var Methods = [...]string{
	http.MethodConnect, http.MethodDelete, http.MethodGet,
	http.MethodHead, http.MethodOptions, http.MethodPatch,
	http.MethodPost, http.MethodPut, http.MethodTrace,
}

func HasRequestTitle(payload []byte) bool {
	s := SliceToString(payload)
	if len(s) < MinRequestCount {
		return false
	}
	titleLen := bytes.Index(payload, CRLF)
	if titleLen == -1 {
		return false
	}
	if strings.Count(s[:titleLen], " ") != 2 {
		return false
	}
	method := string(Method(payload))
	var methodFound bool
	for _, m := range Methods {
		if methodFound = method == m; methodFound {
			break
		}
	}
	if !methodFound {
		return false
	}
	path := strings.Index(s[len(method)+1:], " ")
	if path == -1 {
		return false
	}
	major, minor, ok := http.ParseHTTPVersion(s[path+len(method)+2 : titleLen])
	return ok && major == 1 && (minor == 0 || minor == 1)
}

func HasResponseTitle(payload []byte) bool {
	s := SliceToString(payload)
	if len(s) < MinResponseCount {
		return false
	}
	titleLen := bytes.Index(payload, CRLF)
	if titleLen == -1 {
		return false
	}
	major, minor, ok := http.ParseHTTPVersion(s[0:VersionLen])
	if !(ok && major == 1 && (minor == 0 || minor == 1)) {
		return false
	}
	if s[VersionLen] != ' ' {
		return false
	}
	status, ok := atoI(payload[VersionLen+1:VersionLen+4], 10)
	if !ok {
		return false
	}
	// only validate status codes mentioned in rfc2616.
	if http.StatusText(status) == "" {
		return false
	}
	// handle cases from #875
	return payload[VersionLen+4] == ' ' || payload[VersionLen+4] == '\r'
}

// SliceToString preferred for large body payload (zero allocation and faster)
func SliceToString(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}

// this works with positive integers
func atoI(s []byte, base int) (num int, ok bool) {
	var v int
	ok = true
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			ok = false
			break
		}
		v = int(hexTable[s[i]])
		if v >= base || (v == 0 && s[i] != '0') {
			ok = false
			break
		}
		num = (num * base) + v
	}
	return
}

var hexTable = [128]byte{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'A': 10,
	'a': 10,
	'B': 11,
	'b': 11,
	'C': 12,
	'c': 12,
	'D': 13,
	'd': 13,
	'E': 14,
	'e': 14,
	'F': 15,
	'f': 15,
}
