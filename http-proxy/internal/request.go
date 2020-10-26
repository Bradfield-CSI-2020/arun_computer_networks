package internal

import (
	"bytes"
	"fmt"
)

type Request struct {
	Status  Status
	Headers [][]byte
	Payload []byte
}

type Status struct {
	data    []byte
	Method  string
	Path    string
	Version string
}

func (req *Request) ToBinary() []byte {
	raw := req.Status.data

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

	proxyHeaders, _ := FilterHopHopHeaders(req.Headers)
	proxyRequest.Headers = proxyHeaders
	proxyRequest.Status = req.Status
	proxyRequest.Payload = req.Payload

	return proxyRequest
}

func (req *Request) ReadRequest(raw []byte) {
	parts := bytes.Split(raw, []byte("\r\n\r\n"))
	nonPayloadParts := bytes.Split(parts[0], []byte("\r\n"))

	req.Payload = parts[1]
	req.setStatus(nonPayloadParts[0])
	req.Headers = nonPayloadParts[1:]
}

func (req *Request) setStatus(raw []byte) {
	statusString := bytes.Split(raw, []byte(" "))
	req.Status.data = raw
	req.Status.Method = string(statusString[0])
	req.Status.Path = string(statusString[1])
	req.Status.Version = string(statusString[2])
}

func (req *Request) Print() {

	fmt.Println("----Status----")
	fmt.Printf("Method    : %s\n", req.Status.Method)
	fmt.Printf("Path      : %s\n", req.Status.Path)
	fmt.Printf("Version   : %s\n", req.Status.Version)

	fmt.Printf("Request headers   : %s\n", req.Headers)
	fmt.Printf("Request payload   : %s\n", req.Payload)
	fmt.Printf("Payload length    : %d\n", len(req.Payload))
	fmt.Println("----")
}
