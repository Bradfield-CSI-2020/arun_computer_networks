package internal

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

func FilterHopHopHeaders(raw [][]byte) ([][]byte, [][]byte)  {

	HopHopHeaders := [][]byte{[]byte(Connection), []byte(KeepAlive), []byte(TransferEncoding), []byte(TE), []byte(Trailer), []byte(Upgrade), []byte(ProxyAuthorization), []byte(ProxyAuthenticate)}
	var carryForwardHeaders [][]byte
	var hopHeaders [][]byte

	for _, headerBytes := range raw {
		toFilter := false

		for _, hopHeader := range HopHopHeaders {
			if bytes.HasPrefix(headerBytes, hopHeader) {
				toFilter = true
				hopHeaders  =  append(hopHeaders, hopHeader)
				continue
			}
		}

		if !toFilter {
			carryForwardHeaders = append(carryForwardHeaders, headerBytes)
		}
	}

	return carryForwardHeaders, hopHeaders
}
