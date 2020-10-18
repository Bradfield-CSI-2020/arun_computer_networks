package parse

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
)

var data []byte
var byteOrder binary.ByteOrder

func init() {

	fileBytes, readError := ioutil.ReadFile("net.cap")

	if readError != nil {
		log.Fatalf("parse:readError %v\n", readError)
	}

	data = fileBytes
	mgBuffer := bytes.NewBuffer([]byte{fileBytes[0], fileBytes[1], fileBytes[2], fileBytes[3]})

	var magicNumber uint32

	err := binary.Read(mgBuffer, binary.LittleEndian, &magicNumber)

	if err != nil {
		log.Fatalf("parse:magicNumber %v\n", err)
	}

	// TODO: why does this return a different endianness?
	// rep := hex.EncodeToString(magicNumberBytes)
	// fmt.Printf("hex: %s\n", rep)

	// TODO: what to do about this:
	// with the two nibbles of the two lower-order bytes of the magic number swapped
	if magicNumber == 0xa1b2c3d4 {
		byteOrder = binary.LittleEndian
		fmt.Println("host byte order: LittleEndian")
	} else if magicNumber == 0xd4c3b2a1 {
		byteOrder = binary.BigEndian
		fmt.Println("host byte order: BigEndian")
	} else {
		log.Fatalf("unknown byte order found: %x", magicNumber)
	}

}

func PcapHeader() {

	majorVersionRaw := []byte{data[4],data[5]}
	minorVersionRaw := []byte{data[6],data[7]}

	var majorVersion uint16
	var minorVersion uint16

	mgvBuffer := bytes.NewBuffer(majorVersionRaw)
	err := binary.Read(mgvBuffer, byteOrder, &majorVersion)

	if err != nil {
		log.Fatalf("parse:magicNumber %v\n", err)
	}

	mivBuffer := bytes.NewBuffer(minorVersionRaw)
	err = binary.Read(mivBuffer, byteOrder, &minorVersion)

	if err != nil {
		log.Fatalf("parse:magicNumber %v\n", err)
	}

	fmt.Printf("version is %d.%d\n",majorVersion, minorVersion)

}
