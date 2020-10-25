package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	service := ":3001"
	tcpAddrPointer, err := net.ResolveTCPAddr("tcp", service)
	assertNil(err, "")
	listener, err := net.ListenTCP("tcp", tcpAddrPointer)
	assertNil(err, "")

	fmt.Println("starting tcp server...")

	for {
		conn, err := listener.Accept()
		assertNil(err, "")
		fmt.Println("receiving a request...")

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		assertNil(err, "")

		fmt.Println("received message: ", string(buf))

		_, err = conn.Write(buf)
		assertNil(err, "")

		err = conn.Close()
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
