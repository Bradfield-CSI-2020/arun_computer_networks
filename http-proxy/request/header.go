package request

import (
	"bytes"
	"fmt"
)

type Request struct {
	Status  []byte
	Headers [][]byte
	Payload []byte
}

func (req *Request) ToBinary() []byte {
	raw := req.Status

	for _, v := range req.Headers {
		raw = append(raw, []byte("\r\n")...)
		raw = append(raw, v...)
	}

	raw = append(raw, []byte("\r\n\r\n")...)
	raw = append(raw, req.Payload...)

	return raw
}

func (req *Request) GenerateProxyRequest() Request {
	var proxyRequest Request

	proxyHeaders := FilterHopHopHeaders(req.Headers)
	proxyRequest.Headers = proxyHeaders
	proxyRequest.Status = req.Status
	proxyRequest.Payload = req.Payload

	return proxyRequest
}

func (req *Request) ReadRequest(raw []byte) {
	parts := bytes.Split(raw, []byte("\r\n\r\n"))
	nonPayloadPart := bytes.Split(parts[0], []byte("\r\n"))

	req.Payload = parts[1]
	req.Status = nonPayloadPart[0]
	req.Headers = nonPayloadPart[1:]
}

func (req *Request) Print() {
	fmt.Println("----")
	fmt.Printf("Request status	: %s\n", req.Status)
	fmt.Printf("Request headers	: %s\n", req.Headers)
	fmt.Printf("Request payload	: %s\n", req.Payload)
	fmt.Printf("Payload length	: %d\n", len(req.Payload))
	fmt.Println("----")
}
