package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	service := ":3001"

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)

	if err != nil {
		log.Fatalf("error resolving tcp address: %v\n ", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Fatalf("error resolving tcp address: %v\n ", err)
	}

	fmt.Println("starting tcp server...")

	for {
		conn, err := listener.Accept()

		fmt.Println("receiving a request...")

		if err != nil {
			log.Fatalf("error resolving tcp address: %v\n ", err)
		}

		daytime := time.Now().String()

		_, err = conn.Write([]byte(daytime))

		if err != nil {
			log.Fatalf("error writing to conn: %v\n ", err)
		}

		err = conn.Close()

		if err != nil {
			log.Fatalf("error closing conn: %v\n ", err)
		}
	}
}