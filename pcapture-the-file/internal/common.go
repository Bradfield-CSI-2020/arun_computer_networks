package internal

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Parse4ByteValue(raw []byte, byteOrder binary.ByteOrder) uint32 {


	if len(raw) != 4 {
		log.Fatal("Parse4ByteValue: raw is not the right size")
	}

	buffer := bytes.NewBuffer(raw)

	var container uint32

	err := binary.Read(buffer, byteOrder, &container)

	if err != nil {
		log.Fatalf("parse:Parse4ByteValue %v\n", err)
	}

	return container
}

func Parse2ByteValue(raw []byte, byteOrder binary.ByteOrder) uint16 {
	if len(raw) != 2 {
		log.Fatal("Parse2ByteValue: raw is not the right size")
	}

	buffer := bytes.NewBuffer(raw)

	var container uint16

	err := binary.Read(buffer, byteOrder, &container)

	if err != nil {
		log.Fatalf("parse:Parse2ByteValue %v\n", err)
	}
	return container
}
