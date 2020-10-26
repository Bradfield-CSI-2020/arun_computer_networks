package main

import (
	"fmt"
	"http_proxy/cache"
	"http_proxy/request"
	"io/ioutil"
	"log"
	"net"
)

func main() {

	//setup cache
	proxyCache := cache.InitCache()

	// setup proxy server
	proxyServerAddr, err := net.ResolveTCPAddr("tcp", "localhost:3001")
	assertNil(err, "")

	proxyListener, err := net.ListenTCP("tcp", proxyServerAddr)
	assertNil(err, "")

	fmt.Println("starting proxy server...")

	for {
		proxyConn, err := proxyListener.Accept()
		assertNil(err, "")

		buf := make([]byte, 2048)
		n, err := proxyConn.Read(buf)
		assertNil(err, "")

		var incomingRequest request.Request
		incomingRequest.ReadRequest(buf[0:n])

		cachedValue := proxyCache.Get(incomingRequest.Status.Path)

		targetServerAddr, err := net.ResolveTCPAddr("tcp", "localhost:9000")
		assertNil(err, "")

		serverConn, err := net.DialTCP("tcp", nil, targetServerAddr)
		assertNil(err, "")

		if cachedValue != nil {
			fmt.Println("returning value from cache")
			_, err = proxyConn.Write(cachedValue)
		} else {
			check := incomingRequest.ToBinary()
			fmt.Println("original size: ", len(buf[0:n]))
			fmt.Println("check size: ", len(check))
			incomingRequest.Print()

			var proxyRequest = incomingRequest.GenerateProxyRequest()
			binary := proxyRequest.ToBinary()
			_, err = serverConn.Write(binary)
			assertNil(err, "")

			result, err := ioutil.ReadAll(serverConn)
			assertNil(err, "")

			proxyCache.Set(incomingRequest.Status.Path, result)

			_, err = proxyConn.Write(result)
			assertNil(err, "")
		}

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
