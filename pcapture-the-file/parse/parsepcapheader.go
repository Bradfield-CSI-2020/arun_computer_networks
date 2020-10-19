package parse

import (
	"computer_networks/internal"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
)

var Data []byte
var HostByteOrder binary.ByteOrder
var RawPacketData []byte

func init() {

	fileBytes, readError := ioutil.ReadFile("net.cap")

	if readError != nil {
		log.Fatalf("parse:readError %v\n", readError)
	}

	Data = fileBytes

	fmt.Printf("file size is %d bytes\n", len(Data))

	magicNumber := internal.Parse4ByteValue(fileBytes[0:4], binary.LittleEndian)

	// TODO: what to do about this:
	// with the two nibbles of the two lower-order bytes of the magic number swapped
	if magicNumber == 0xa1b2c3d4 {
		HostByteOrder = binary.LittleEndian
	} else if magicNumber == 0xd4c3b2a1 {
		HostByteOrder = binary.BigEndian
	} else {
		log.Fatalf("unknown byte order found: %x", magicNumber)
	}

	RawPacketData = fileBytes[24:]

	fmt.Printf("packet data size is %d bytes\n", len(RawPacketData))
}

type PCapFileHeader struct {
	MajorVersion uint16
	MinorVersion uint16
	TimestampOffset uint32
	TimestampAccuracy uint32
	SnapShotLength uint32
	LinkLayerHeaderType uint32
}

type PcapPacketHeader struct {
	TimestampSeconds uint32
	TimestampMicroSeconds  uint32
	PacketLength uint32
	FullPacketLength uint32
}

type PacketData struct {

}

func GetPacketData() {

	startIndex := 0
	count := 0

	for  startIndex < len(RawPacketData) {
		headerData := RawPacketData[startIndex:startIndex+16]

		parsedData := PacketHeader(headerData)

		if parsedData.PacketLength != 78 {
			fmt.Printf("packet length: %d\n", parsedData.PacketLength)
			fmt.Printf("full packet length: %d\n", parsedData.FullPacketLength)
		}

		if parsedData.PacketLength != parsedData.PacketLength {
			log.Fatalf("packet length mismatch")
		}

		startIndex = startIndex + int(parsedData.PacketLength) + 16
		count = count + 1

	}

	fmt.Printf("count: %d\n", count)
}

func PacketHeader(raw []byte) PcapPacketHeader {
	var header PcapPacketHeader

	timestampData := raw[0:4]
	timestampNanoData := raw[4:8]
	packetLengthData := raw[8:12]
	fullPacketLengthData := raw[12:16]

	header.TimestampSeconds = internal.Parse4ByteValue(timestampData, HostByteOrder)
	header.TimestampMicroSeconds = internal.Parse4ByteValue(timestampNanoData, HostByteOrder)
	header.PacketLength = internal.Parse4ByteValue(packetLengthData, HostByteOrder)
	header.FullPacketLength = internal.Parse4ByteValue(fullPacketLengthData, HostByteOrder)

	return header
}

func PcapHeader() PCapFileHeader {
	var header PCapFileHeader
	x, y := parseVersion(Data)

	header.MajorVersion = x
	header.MinorVersion = y
	header.TimestampOffset = parseTimeStampOffset(Data)
	header.TimestampAccuracy = parseTimeStampAccuracy(Data)
	header.SnapShotLength = parseSnapShotLength(Data)
	header.LinkLayerHeaderType = parseLinkLayerHeaderType(Data)

	return header
}

func parseVersion(pCapHeader [] byte) (x uint16, y uint16) {
	majorVersionRaw := pCapHeader[4:6]
	x = internal.Parse2ByteValue(majorVersionRaw, HostByteOrder)

	minorVersionRaw := pCapHeader[6:8]
	y = internal.Parse2ByteValue(minorVersionRaw, HostByteOrder)

	return x,y
}

func parseTimeStampOffset(pCapHeader [] byte) uint32{
	raw := pCapHeader[8:12]
	return internal.Parse4ByteValue(raw, HostByteOrder)
}

func parseTimeStampAccuracy(pCapHeader [] byte) uint32 {
	timeZoneOffsetRaw := pCapHeader[12:16]
	return internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder)
}

func parseSnapShotLength(pCapHeader [] byte) uint32 {
	timeZoneOffsetRaw := pCapHeader[16:20]
	return internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder)
}

func parseLinkLayerHeaderType(pCapHeader [] byte) uint32 {
	timeZoneOffsetRaw := pCapHeader[20:24]
	return internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder)
}
