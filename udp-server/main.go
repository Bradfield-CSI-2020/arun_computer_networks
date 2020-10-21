package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp4", ":3002")

	if err != nil {
		log.Fatalf("error resolving udp address: %v\n ", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatalf("error listening udp address: %v\n ", err)
	}

	fmt.Println("starting udp server...")

	for {
		handleClient(conn)
	}


}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte

	_, addr, err := conn.ReadFromUDP(buf[0:])

	fmt.Println("received udp request...")

	if err != nil {
		log.Fatalf("error reading from udp address: %v\n ", err)
	}

	daytime := time.Now().String()
	_, err = conn.WriteToUDP([]byte(daytime), addr)

	if err != nil {
		log.Fatalf("error writing to udp address: %v\n ", err)
	}
}
