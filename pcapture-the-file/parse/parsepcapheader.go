package parse

import (
	"computer_networks/internal"
	"encoding/binary"
	"log"
)

type MyFile struct {
	PCapHeader     PCapFileHeader
	PacketDataData []PacketData
	IpDataGramData []IpDataGram
}

type PCapFileHeader struct {
	NewHostByteOrder    binary.ByteOrder
	MajorVersion        uint16
	MinorVersion        uint16
	TimestampOffset     uint32
	TimestampAccuracy   uint32
	SnapShotLength      uint32
	LinkLayerHeaderType uint32
}

type PacketData struct {
	PacketHeader  PcapPacketHeader
	RawData       []byte
	EtherNetFrame EtherNetFrameRaw
}

type PcapPacketHeader struct {
	TimestampSeconds      uint32
	TimestampMicroSeconds uint32
	PacketLength          uint32
	FullPacketLength      uint32
}

type EtherNetFrameRaw struct {
	DestinationMac []byte
	SourceMac      []byte
	EtherType      []byte
	IpRawPayload   []byte
	InterPacketGap []byte
}

type IpHeader struct {
	InternetHeaderLength uint32
	TotalLength          uint32
	ECN			         uint32
	Protocol             uint32
	SourceIp             uint32
	DestinationIp        uint32
}

type IpDataGram struct {
	IpHeader IpHeader
	RawData  []byte
}

func ReadIpHeader(raw []byte) IpHeader {
	var header IpHeader

	header.InternetHeaderLength = uint32(raw[0]&0b1111) * 4
	header.TotalLength = uint32(internal.Parse2ByteValueBig(raw[2:4]))
	header.ECN = uint32(raw[1])&0x3
	header.Protocol = uint32(raw[9])
	header.SourceIp = internal.Parse4ByteValueBig(raw[12:16])
	header.DestinationIp = internal.Parse4ByteValueBig(raw[16:20])

	return header
}

func ReadEtherNetHeaders(raw []byte) EtherNetFrameRaw {

	var etherNetFrame EtherNetFrameRaw

	etherNetFrame.DestinationMac = raw[0:6]
	etherNetFrame.SourceMac = raw[6:12]
	etherNetFrame.EtherType = raw[12:14]
	etherNetFrame.IpRawPayload = raw[14 : len(raw)-12]
	etherNetFrame.InterPacketGap = raw[len(raw)-12:]

	return etherNetFrame
}

func ReadPacketHeader(raw []byte) PcapPacketHeader {
	var header PcapPacketHeader

	timestampData := raw[0:4]
	timestampNanoData := raw[4:8]
	packetLengthData := raw[8:12]
	fullPacketLengthData := raw[12:16]

	header.TimestampSeconds = internal.Parse4ByteValue(timestampData)
	header.TimestampMicroSeconds = internal.Parse4ByteValue(timestampNanoData)
	header.PacketLength = internal.Parse4ByteValue(packetLengthData)
	header.FullPacketLength = internal.Parse4ByteValue(fullPacketLengthData)

	return header
}

func ReadPCapHeader(raw []byte) PCapFileHeader {

	var header PCapFileHeader
	magicNumber := internal.Parse4ByteValue(raw[0:4])

	if magicNumber == 0xa1b2c3d4 {
		header.NewHostByteOrder = binary.LittleEndian
	} else if magicNumber == 0xd4c3b2a1 {
		header.NewHostByteOrder = binary.BigEndian
	} else {
		log.Fatalf("unknown byte order found: %x", magicNumber)
	}

	x, y := parseVersion(raw)

	header.MajorVersion = x
	header.MinorVersion = y
	header.TimestampOffset = parseTimeStampOffset(raw)
	header.TimestampAccuracy = parseTimeStampAccuracy(raw)
	header.SnapShotLength = parseSnapShotLength(raw)
	header.LinkLayerHeaderType = parseLinkLayerHeaderType(raw)

	return header
}

func parseVersion(pCapHeader []byte) (x uint16, y uint16) {
	majorVersionRaw := pCapHeader[4:6]
	x = internal.Parse2ByteValue(majorVersionRaw)

	minorVersionRaw := pCapHeader[6:8]
	y = internal.Parse2ByteValue(minorVersionRaw)

	return x, y
}

func parseTimeStampOffset(pCapHeader []byte) uint32 {
	raw := pCapHeader[8:12]
	return internal.Parse4ByteValue(raw)
}

func parseTimeStampAccuracy(pCapHeader []byte) uint32 {
	timeZoneOffsetRaw := pCapHeader[12:16]
	return internal.Parse4ByteValue(timeZoneOffsetRaw)
}

func parseSnapShotLength(pCapHeader []byte) uint32 {
	timeZoneOffsetRaw := pCapHeader[16:20]
	return internal.Parse4ByteValue(timeZoneOffsetRaw)
}

func parseLinkLayerHeaderType(pCapHeader []byte) uint32 {
	timeZoneOffsetRaw := pCapHeader[20:24]
	return internal.Parse4ByteValue(timeZoneOffsetRaw)
}
