package main

import (
	"dns_client/dns_message"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("expected 1 args -> got %d args instead", len(os.Args)-1)
	}

	makeDnsRequest(os.Args[1])
}

func makeDnsRequest(domainName string) {

	udpAddr, err := net.ResolveUDPAddr("udp4", "8.8.8.8:53")

	if err != nil {
		log.Fatalf("error resolving udp address: %v\n ", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		log.Fatalf("error dialing udp: %v\n ", err)
	}

	var query dns_message.DnsQuery

	query.InitQuery(domainName)

	query.Print()

	payload := query.GenerateBinaryPayload()

	fmt.Println("size of query: ", len(payload))

	_, err = conn.Write(payload)

	if err != nil {
		log.Fatalf("error writing to udp socket: %v\n ", err)
	}

	var buf [512]byte

	n, err := conn.Read(buf[0:])

	if err != nil {
		log.Fatalf("error reading from udp socket: %v\n ", err)
	}

	fmt.Println("size of response: ", n)

	var reply dns_message.DnsReply

	reply.ReadPayload(buf[0:n], domainName)

	reply.Print()

	err = conn.Close()

	if err != nil {
		log.Fatalf("error closing udp socket: %v\n ", err)
	}

	os.Exit(0)
}
