package main

import (
	"fmt"
	"http_proxy/request"
	"io/ioutil"
	"log"
	"net"
)

func main() {

	// setup server
	proxyServerAddr, err := net.ResolveTCPAddr("tcp", "localhost:3001")
	assertNil(err, "")

	proxyListener, err := net.ListenTCP("tcp", proxyServerAddr)
	assertNil(err, "")

	fmt.Println("starting tcp  ...")

	for {
		proxyConn, err := proxyListener.Accept()
		assertNil(err, "")

		buf := make([]byte, 2048)
		n, err := proxyConn.Read(buf)
		assertNil(err, "")

		var request internal.Request
		request.ReadRequest(buf[0:n])

		request.Print()

		var proxyRequest = request.GenerateProxyRequest()

		binary := proxyRequest.ToBinary()

		targetServerAddr, err := net.ResolveTCPAddr("tcp", "localhost:9000")
		assertNil(err, "")

		serverConn, err := net.DialTCP("tcp", nil, targetServerAddr)
		assertNil(err, "")

		_, err = serverConn.Write(binary)
		assertNil(err, "")

		result, err := ioutil.ReadAll(serverConn)
		assertNil(err, "")

		fmt.Printf("received message from target: %s\n", result)

		_, err = proxyConn.Write(result)
		assertNil(err, "")

		err = serverConn.Close()
		assertNil(err, "")

		err = proxyConn.Close()
		assertNil(err, "")

	}
}

func assertNil(err error, message string) {
	if message == "" {
		message = "something bad happened: %v\n"
	}
	if err != nil {
		log.Fatalf(message, err)
	}
}
