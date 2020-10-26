package request

import (
	"bytes"
)

const (
	Connection         = "Connection:"
	KeepAlive          = "Keep-Alive:"
	TransferEncoding   = "Transfer-Encoding:"
	TE                 = "TE:"
	Trailer            = "Trailer:"
	Upgrade            = "Upgrade:"
	ProxyAuthorization = "Proxy-Authorization:"
	ProxyAuthenticate  = "Proxy-Authenticate:"
)

func FilterHopHopHeaders(raw [][]byte) [][]byte {

	HopHopHeaders := [][]byte{[]byte(Connection), []byte(KeepAlive), []byte(TransferEncoding), []byte(TE), []byte(Trailer), []byte(Upgrade), []byte(ProxyAuthorization), []byte(ProxyAuthenticate)}
	var filteredHeaders [][]byte

	for _, headerBytes := range raw {
		toFilter := false

		for _, hopHeader := range HopHopHeaders {
			if bytes.HasPrefix(headerBytes, hopHeader) {
				toFilter = true
				continue
			}
		}

		if !toFilter {
			filteredHeaders = append(filteredHeaders, headerBytes)
		}
	}

	return filteredHeaders
}
