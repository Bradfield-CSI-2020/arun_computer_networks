package cmd

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

	var magicNumber uint32

	internal.Parse4ByteValue([]byte{fileBytes[0], fileBytes[1], fileBytes[2], fileBytes[3]}, binary.LittleEndian,  &magicNumber)

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

func PcapHeader() {
	parseVersion()
	parseTimeStampOffset()
	parseTimeStampAccuracy()
	parseSnapShotLength()
	parseLinkLayerHeaderType()
}

func parseVersion() {
	majorVersionRaw := []byte{Data[4], Data[5]}
	minorVersionRaw := []byte{Data[6], Data[7]}

	var majorVersion uint16
	var minorVersion uint16

	internal.Parse2ByteValue(majorVersionRaw, HostByteOrder, &majorVersion)
	internal.Parse2ByteValue(minorVersionRaw, HostByteOrder, &minorVersion)

	fmt.Printf("version is %d.%d\n",majorVersion, minorVersion)
}

func parseTimeStampOffset() {
	raw := []byte{Data[8], Data[9], Data[10], Data[11]}

	var value uint32

	internal.Parse4ByteValue(raw, HostByteOrder, &value)

	fmt.Printf("timestamp offset is %d\n", value)
}

func parseTimeStampAccuracy() {
	timeZoneOffsetRaw := []byte{Data[12], Data[13], Data[14], Data[15]}

	var majorVersion uint32

	internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder, &majorVersion)

	fmt.Printf("timestamp accuracy is %d\n", majorVersion)
}

func parseSnapShotLength() {
	timeZoneOffsetRaw := []byte{Data[16], Data[17], Data[18], Data[19]}

	var majorVersion uint32

	internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder, &majorVersion)

	fmt.Printf("snap shot length is %d\n", majorVersion)
}

func parseLinkLayerHeaderType() {
	timeZoneOffsetRaw := []byte{Data[20], Data[21], Data[22], Data[23]}

	var majorVersion uint32

	internal.Parse4ByteValue(timeZoneOffsetRaw, HostByteOrder, &majorVersion)

	fmt.Printf("link layer header type %d\n", majorVersion)
}