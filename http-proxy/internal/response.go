package internal

import (
	"bytes"
	"fmt"
)

type Response struct {
	Status  ResponseStatus
	Headers [][]byte
	Payload []byte
}

type ResponseStatus struct {
	data         []byte
	Version      string
	StatusCode   string
	ShortMessage string
}

func (res *Response) ReadResponse(raw []byte) {
	parts := bytes.Split(raw, []byte("\r\n\r\n"))

	if len(parts) != 2 {
		fmt.Printf("Expected size 2 but got %d !!\n", len(parts))
	}

	for i, v := range parts {
		fmt.Printf("parts %d:", i)
		fmt.Printf("%s\n", v)
	}

	nonPayloadParts := bytes.Split(parts[0], []byte("\r\n"))

	res.Payload = parts[1]
	res.setStatus(nonPayloadParts[0])
	res.Headers = nonPayloadParts[1:]
}

func (res *Response) setStatus(raw []byte) {
	statusString := bytes.Split(raw, []byte(" "))
	res.Status.data = raw
	res.Status.Version = string(statusString[0])
	res.Status.StatusCode = string(statusString[1])
	res.Status.ShortMessage = string(statusString[2])
}

func (res *Response) Print() {

	fmt.Println("----Response Status----")
	fmt.Printf("Version    : %s\n", res.Status.Version)
	fmt.Printf("StatusCode      : %s\n", res.Status.StatusCode)
	fmt.Printf("ShortMessage   : %s\n", res.Status.ShortMessage)

	fmt.Printf("Response headers   : %s\n", res.Headers)
	fmt.Printf("Response payload   : %s\n", res.Payload)
	fmt.Printf("Payload length    : %d\n", len(res.Payload))
	fmt.Println("-----------END---------")
}
