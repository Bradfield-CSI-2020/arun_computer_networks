package main

import (
	"dns_client/dns_message"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	MakeUdpRequest()
}

func MakeUdpRequest() {
	udpAddr, err := net.ResolveUDPAddr("udp4", "8.8.8.8:53")

	if err != nil {
		log.Fatalf("error resolving udp address: %v\n ", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		log.Fatalf("error dialing udp: %v\n ", err)
	}

	message := dns_message.InitQuery("google.com")

	messageBin := message.GenerateBinaryPayload()

	_, err = conn.Write(messageBin)

	if err != nil {
		log.Fatalf("error writing to udp socket: %v\n ", err)
	}

	var buf [512]byte

	n, err := conn.Read(buf[0:])

	if err != nil {
		log.Fatalf("error reading from udp socket: %v\n ", err)
	}

	fmt.Println("size of response: ", n)

	response := dns_message.ReadPayload(buf[0:n])

	response.Print()

	err = conn.Close()

	if err != nil {
		log.Fatalf("error closing udp socket: %v\n ", err)
	}

	os.Exit(0)
}

func MakeTcpRequest() {
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
