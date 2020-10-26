package dns_message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Base Message

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
	QR     uint8 // 0 or 1 (1 bit) - 0 for QUERY, 1 for RESPONSE
	OPCode uint8 // 0000 -> for standard message QUERY
	AA     uint8
	TC     uint8
	RD     uint8
	RA     uint8
	z      uint8 // 0000 -> 4 bits all zero
	RCode  uint8 // 0000 -> 4 bits max
}

type nameServerQuestion struct {
	Name   string
	QType  uint16 // set to standard?
	QClass uint16 // se to IN?
}

type answer struct {
	Name     []byte
	Type     []byte
	Class    []byte
	TTL      uint32
	RDLength uint16
	RData    []byte
}

type authority struct {
}

type additional struct {
}

// Query
type DnsQuery struct {
	message dnsMessage
}

func (query DnsQuery) InitQuery(domainName string) DnsQuery {

	var message dnsMessage

	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)

	message.header.Id = uint16(generator.Int31n(1 << 16)) // TODO: make this random
	message.header.Flags.QR = 0
	message.header.Flags.OPCode = 0
	message.header.QDCount = 1

	// NOT required for Query
	message.header.ANCount = 0
	message.header.NSCount = 0
	message.header.ARCount = 0

	message.question.Name = domainName

	query.message = message
	return query
}

func (query *DnsQuery) GenerateBinaryPayload() []byte {

	message := query.message
	// GenerateBinaryPayload Header
	id := make([]byte, 2)

	binary.BigEndian.PutUint16(id, message.header.Id)

	flags := make([]byte, 2)

	flags[0] = byte(0) | // init zero byte
		(message.header.Flags.QR << 7) | // set QR     // TODO: no need to shift just set to zero
		byte(0) | // set OPCODE (NOT NEEDED) all zero for QUERY
		byte(0) | // set AA (NOT NEEDED) all zero
		byte(0) | //
		(byte(1) << 1) | // set RD as required 	0000 0010
		byte(0) // set RA to zero

	flags[1] = byte(0) // Z is all zero , RCODE is all zero for a internal

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

func (query *DnsQuery) Print() {
	message := query.message

	fmt.Printf("-------Query Headers -------:\n")
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

	fmt.Println()
	fmt.Printf("-------Query Question -------:\n")

	fmt.Printf("QName: %d\n", message.question.Name)
	fmt.Printf("QClass: %d\n", message.question.QClass)
	fmt.Printf("QType: %d\n", message.question.QType)
}

// DNS Reply
type DnsReply struct {
	message dnsMessage
}

func (reply DnsReply) ReadPayload(raw []byte, domainName string) DnsReply {

	var message dnsMessage
	headerParts := raw[0:12]

	message.header.Id = binary.BigEndian.Uint16(raw[0:2])

	message.header.Id = binary.BigEndian.Uint16(raw[0:2])

	flagsRaw := raw[2:4]

	// First 8 bits
	message.header.Flags.QR = (flagsRaw[0] & (uint8(255))) >> 7
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

	encodedDomainName := encodeDomainName(domainName)
	questionRaw := encodedDomainName

	qType := []byte{uint8(0), uint8(1)}
	qClass := []byte{uint8(0), uint8(1)}

	// append headers and question payloads
	questionRaw = append(questionRaw, qType...)
	questionRaw = append(questionRaw, qClass...)

	parts := bytes.Split(raw, questionRaw)

	answerPart := parts[1]

	message.answer.Name = answerPart[0:2] // name is compressed and so is 2 octets
	message.answer.Type = answerPart[2:4]
	message.answer.Class = answerPart[4:6]
	message.answer.TTL = binary.BigEndian.Uint32(answerPart[6:10])       // seconds
	message.answer.RDLength = binary.BigEndian.Uint16(answerPart[10:12]) // octet count
	message.answer.RData = answerPart[12 : 12+message.answer.RDLength]

	reply.message = message
	return reply

}

func (reply *DnsReply) Print() {

	message := reply.message
	fmt.Println()
	fmt.Printf("-------Reply Headers -------:\n")
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

	fmt.Println()
	fmt.Printf("-------Reply Ansswer -------:\n")
	fmt.Printf("Name: %d\n", message.answer.Name)
	fmt.Printf("Type: %d\n", message.answer.Type)
	fmt.Printf("Class: %d\n", message.answer.Class)
	fmt.Printf("TTL: %d\n", message.answer.TTL)
	fmt.Printf("RD Length: %d\n", message.answer.RDLength)

	var ipParts []string

	for _, v := range message.answer.RData {
		ipParts = append(ipParts, strconv.Itoa(int(v)))
	}

	fmt.Println("IP Address: ", strings.Join(ipParts, "."))

}

// encode domain name labels into question format
func encodeDomainName(domainName string) []byte {
	labels := strings.Split(domainName, ".")
	var raw []byte
	for _, label := range labels {
		size := uint8(len(label))
		raw = append(raw, size)
		raw = append(raw, []byte(label)...)
	}

	raw = append(raw, byte(0))
	return raw
}
