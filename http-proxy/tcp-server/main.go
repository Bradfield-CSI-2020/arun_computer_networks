package main

import (
	"fmt"
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

	fmt.Println("starting tcp server...")

	for {
		proxyConn, err := proxyListener.Accept()
		assertNil(err, "")

		buf := make([]byte, 2048)
		_, err = proxyConn.Read(buf)
		assertNil(err, "")

		fmt.Println("received message: \n", string(buf))

		targetServerAddr, err := net.ResolveTCPAddr("tcp", "localhost:9000")
		assertNil(err, "")

		serverConn, err := net.DialTCP("tcp", nil, targetServerAddr)
		assertNil(err, "")

		_, err = serverConn.Write(buf)

		//_, err = serverConn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
		//_, err = serverConn.Write([]byte("GET / HTTP/1.0\r\nHost: www.arun.com\r\nConnection: close\r\nUser-Agent: arun\r\n\r\n"))
		result, err := ioutil.ReadAll(serverConn)

		fmt.Println("received message from target: ", string(result))
		assertNil(err, "")

		assertNil(err, "")
		err = serverConn.Close()
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
