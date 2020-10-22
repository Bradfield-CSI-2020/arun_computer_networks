package dns_message

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type dnsMessage struct {
	header     messageHeader
	question   nameServerQuestion
	answer     answer
	authority  authority
	additional additional
}

// 12 bytes - 6 * (2 byte) sections
type messageHeader struct {
	Id      uint16
	Flags   messageHeaderFlags
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type messageHeaderFlags struct {
	QR     byte // 0 or 1 (1 bit) - 0 for QUERY, 1 for RESPONSE
	OPCode byte // 0000 -> for standard message QUERY
	AA     byte
	TC     byte
	RD     byte
	RA     byte
	z      byte // 0000 -> 4 bits all zero
	RCode  byte // 0000 -> 4 bits max
}

type nameServerQuestion struct {
	Name   string
	QType  uint16 // set to standard?
	QClass uint16 // se to IN?
}

type answer struct {
}

type authority struct {
}

type additional struct {
}

func InitQuery(domainName string) dnsMessage {

	var message dnsMessage

	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)

	message.header.Id = uint16(generator.Int31n(1 << 16))  // TODO: make this random
	message.header.Flags.QR = 0
	message.header.Flags.OPCode = 0
	message.header.QDCount = 1

	// NOT required for Query
	message.header.ANCount = 0
	message.header.NSCount = 0
	message.header.ARCount = 0

	message.question.Name = domainName

	return message
}

func (message *dnsMessage) Print() {
	fmt.Printf("Reply ID: %d\n", message.header.Id)
	fmt.Printf("QR Flag: %d\n", message.header.Flags.QR)
	fmt.Printf("OPCode Flag: %d\n", message.header.Flags.OPCode)
	fmt.Printf("AA Flag: %d\n", message.header.Flags.AA)
	fmt.Printf("TC Flag: %d\n", message.header.Flags.TC)
	fmt.Printf("RD Flag: %d\n", message.header.Flags.RD)
	fmt.Printf("RA Flag: %d\n", message.header.Flags.RA)
	fmt.Printf("RCode Flag: %d\n", message.header.Flags.RCode)
	fmt.Printf("No. of Questions: %d\n", message.header.QDCount)
	fmt.Printf("No. of Answers: %d\n", message.header.ANCount)
	fmt.Printf("No. of Name Servers: %d\n", message.header.NSCount)
	fmt.Printf("No. of Authoritative Records: %d\n", message.header.ARCount)
}

func ReadPayload(raw []byte) dnsMessage {

	var message dnsMessage
	headerParts := raw[0:12]

	message.header.Id = binary.BigEndian.Uint16(raw[0:2])

	message.header.Id = binary.BigEndian.Uint16(raw[0:2])

	flagsRaw := raw[2:4]

	// First 8 bits
	message.header.Flags.QR = flagsRaw[0] & byte(1) << 7
	message.header.Flags.OPCode = (flagsRaw[0] & byte(120)) >> 3
	message.header.Flags.AA = (flagsRaw[0] & (byte(1) << 2)) >> 2
	message.header.Flags.TC = (flagsRaw[0] & (byte(1) << 1)) >> 1
	message.header.Flags.RD = flagsRaw[0] & byte(1)


	// Second 8 bits
	// message.header.Flags.Z -> this reserved for future and not needed
	message.header.Flags.RA = (flagsRaw[0] & (byte(1) << 7)) >> 7
	message.header.Flags.RCode = flagsRaw[1] & byte(15)

	message.header.QDCount = binary.BigEndian.Uint16(headerParts[4:6])
	message.header.ANCount = binary.BigEndian.Uint16(headerParts[6:8])
	message.header.NSCount = binary.BigEndian.Uint16(headerParts[8:10])
	message.header.ARCount = binary.BigEndian.Uint16(headerParts[10:12])

	return message

}

func (message *dnsMessage) GenerateBinaryPayload() []byte {

	// GenerateBinaryPayload Header
	id := make([]byte, 2)

	binary.BigEndian.PutUint16(id, message.header.Id)

	flags := make([]byte, 2)

	flags[0] = byte(0) | // init zero byte
		(message.header.Flags.QR << 7) | // set QR     // TODO: no need to shift just set to zero
		byte(0) | // set OPCODE (NOT NEEDED) all zero for QUERY
		byte(0) | // set AA (NOT NEEDED) all zero
		byte(0)	| //
		(byte(1) << 1) | // set RD as required 	0000 0010
		byte(0) // set RA to zero

	flags[1] = byte(0) // Z is all zero , RCODE is all zero for a request

	qCount := make([]byte, 2) // number of questions -> just one for us

	binary.BigEndian.PutUint16(qCount, message.header.QDCount)

	// These are response fields.... zero them out
	ansCount := []byte{0, 0}
	aaCount := []byte{0, 0}
	adtCount := []byte{0, 0}

	// append all message header parts together
	header := append(id, flags...)
	header = append(header, qCount...)
	header = append(header, ansCount...)
	header = append(header, aaCount...)
	header = append(header, adtCount...)

	// generate question payload
	name := encodeDomainName(message.question.Name)
	qType := []byte{uint8(0), uint8(1)}
	qClass := []byte{uint8(0), uint8(1)}

	// append headers and question payloads
	rawMessage := append(header, name...)
	rawMessage = append(rawMessage, qType...)
	rawMessage = append(rawMessage, qClass...)

	return rawMessage

}

// encode domain name labels into question format
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
