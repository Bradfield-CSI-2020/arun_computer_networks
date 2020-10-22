package dns_message

import (
	"encoding/binary"
	"strings"
)

type DnsMessage struct {
	header MessageHeader
	question NameServerQuestion
	//answer Answer
	//authority Authority
	//additional Additional
}

// 12 bytes - 6 * (2 byte) sections
type MessageHeader struct {
	Id uint16
	Flags HeaderFlags
	QDCount uint16
	//ANCount uint16
	//NSCount uint16
	//ARCount uint16
}

type HeaderFlags struct {
	QR byte // 0 or 1 - 1 bit
	OPCode byte // 4 bits
	//AA byte
	//TC byte
	//RD byte
	//RA byte
	//z byte // 0000 -> 4 bits all zero
	//RCode byte /// 0000 - 4 bits max
}

type NameServerQuestion struct {
	Name string
	QType uint16 // set to standard?
	QClass uint16 // se to IN?

}

type Answer struct {

}

type Authority struct {

}

type Additional struct {

}


func (message *DnsMessage) Init(domainName string) *DnsMessage{

	message.header.Id = 256
	message.header.Flags.QR = 0
	message.header.Flags.OPCode = 0
	message.header.QDCount = 1
	//message.header.ANCount = 0
	//message.header.NSCount = 0
	//message.header.ARCount = 0

	message.question.Name = domainName

	return message
}

func (message *DnsMessage) Generate()  []byte {

	// Generate Header
	id := make([]byte,2)

	binary.BigEndian.PutUint16(id, message.header.Id)

	flags := make([]byte,2)

	flags[0] = byte(0) | // init zero byte
		(message.header.Flags.QR << 7) | // set QR
		byte(0) | // set OPCODE (NOT NEEDED) all zero for QUERY
		byte(0) | // set AA (NOT NEEDED) all zero
		(byte(1) << 2) | // set TC as needed? 	0000 0100
		(byte(1) << 1) | // set RD as required 	0000 0010
		byte(0) // set RA to zero

	flags[1] = byte(0) // Z is all zero , RCODE is all zero for a request

	qCount := make([]byte,2)

	binary.BigEndian.PutUint16(qCount, message.header.QDCount)

	ansCount := []byte{0,0}
	aaCount := []byte{0,0}
	adtCount := []byte{0,0}


	header := append(id, flags...)

	header = append(header, qCount...)
	header = append(header, ansCount...)
	header = append(header, aaCount...)
	header = append(header, adtCount...)

	// Payload Section
	name := encodeDomainName(message.question.Name)
	qType := []byte{uint8(0),uint8(1)}
	qClass := []byte{uint8(0),uint8(1)}

	rawMessage := append(header, name...)
	rawMessage = append(rawMessage, qType...)
	rawMessage = append(rawMessage, qClass...)

	return rawMessage

}


func encodeDomainName(domainName string) []byte {
	labels := strings.Split(domainName, ".")
	var bytes []byte
	for _, label := range labels {

		size := uint8(len(label))
		bytes = append(bytes, size)
		bytes = append(bytes, []byte(label)...)
	}

	bytes = append(bytes, byte(0))

	return bytes
}




