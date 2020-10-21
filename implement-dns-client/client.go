package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("hello arun")

	MakeUdpRequest()
}

func MakeUdpRequest() {
	udpAddr, err := net.ResolveUDPAddr("udp4", "localhost:3002")

	if err != nil {
		log.Fatalf("error resolving udp address: %v\n ", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		log.Fatalf("error dialing udp: %v\n ", err)
	}

	_, err = conn.Write([]byte("anything"))

	if err != nil {
		log.Fatalf("error writing to udp socket: %v\n ", err)
	}

	var buf [512]byte

	n, err := conn.Read(buf[0:])

	if err != nil {
		log.Fatalf("error reading from udp socket: %v\n ", err)
	}

	fmt.Println(string(buf[0:n]))

	err = conn.Close()

	if err != nil {
		log.Fatalf("error closing udp socket: %v\n ", err)
	}

	os.Exit(0)
}

func MakeTcpRequest()  {
	fmt.Println("Hello arun !")


	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:3001")

	if err != nil {
		log.Fatalf("error resolving tcp address: %v\n ", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatalf("error dialing tcp: %v\n ", err)
	}

	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))

	result, err := ioutil.ReadAll(conn)

	fmt.Println(string(result))

	if err != nil {
		log.Fatalf("error reading resonse: %v\n ", err)
	}

}