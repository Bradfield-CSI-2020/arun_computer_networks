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

func init() {

	fileBytes, readError := ioutil.ReadFile("net.cap")

	if readError != nil {
		log.Fatalf("parse:readError %v\n", readError)
	}

	Data = fileBytes

	magicRaw := Data[0:4]

	magicNumber := internal.Parse4ByteValue(magicRaw, binary.LittleEndian)


	// TODO: why does this return a different endianness?
	// rep := hex.EncodeToString(magicNumberBytes)
	// fmt.Printf("hex: %s\n", rep)

	// TODO: what to do about this:
	// with the two nibbles of the two lower-order bytes of the magic number swapped
	if magicNumber == 0xa1b2c3d4 {
		HostByteOrder = binary.LittleEndian
		fmt.Println("host byte order: LittleEndian")
	} else if magicNumber == 0xd4c3b2a1 {
		HostByteOrder = binary.BigEndian
		fmt.Println("host byte order: BigEndian")
	} else {
		log.Fatalf("unknown byte order found: %x", magicNumber)
	}

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

func PacketHeader() {

	var header PcapPacketHeader

	raw := Data[24:28]

	header.TimestampSeconds = internal.Parse4ByteValue(raw, HostByteOrder)

	fmt.Printf("timestamp seconds: %d\n", header.TimestampSeconds)

}

func PcapHeader() {
	var header PCapFileHeader
	x, y := parseVersion()

	header.MajorVersion = x
	header.MinorVersion = y
	header.TimestampOffset = parseTimeStampOffset()
	header.TimestampAccuracy = parseTimeStampAccuracy()
	header.SnapShotLength = parseSnapShotLength()
	header.LinkLayerHeaderType = parseLinkLayerHeaderType()

	fmt.Printf("version: %d.%d\n", header.MajorVersion, header.MinorVersion)
	fmt.Printf("timestamp offset: %d\n", header.TimestampOffset)
	fmt.Printf("timestamp accuracy: %d\n", header.TimestampAccuracy)
	fmt.Printf("snapshot length: %d\n", header.SnapShotLength)
	fmt.Printf("linke layer header type: %d\n", header.LinkLayerHeaderType)
}

func parseVersion() (x uint16, y uint16) {
	majorVersionRaw := Data[4:6]
	x = internal.Parse2ByteValue(majorVersionRaw, HostByteOrder)

	minorVersionRaw := []byte{Data[6], Data[7]}
	y = internal.Parse2ByteValue(minorVersionRaw, HostByteOrder)

	return x,y
}

func parseTimeStampOffset() uint32{
	//raw := []byte{Data[8], Data[9], Data[10], Data[11]}
	raw := Data[8:12]
	return internal.Parse4ByteValue(raw, HostByteOrder)
}

func parseTimeStampAccuracy() uint32 {
	//timeZoneOffsetRaw := []byte{Data[12], Data[13], Data[14], Data[15]}
	timeZoneOffsetRaw := Data[12:16]
	return internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder)
}

func parseSnapShotLength() uint32 {
	//timeZoneOffsetRaw := []byte{Data[16], Data[17], Data[18], Data[19]}
	timeZoneOffsetRaw := Data[16:20]
	return internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder)
}

func parseLinkLayerHeaderType() uint32 {
	//timeZoneOffsetRaw := []byte{Data[20], Data[21], Data[22], Data[23]}
	timeZoneOffsetRaw := Data[20:24]
	return internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder)
}
