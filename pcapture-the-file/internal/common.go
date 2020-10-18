package internal

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Parse4ByteValue(raw []byte, byteOrder binary.ByteOrder, container *uint32) {

	if len(raw) != 4 {
		log.Fatal("Parse4ByteValue: raw is not the right size")
	}

	buffer := bytes.NewBuffer(raw)

	err := binary.Read(buffer, byteOrder, container)

	if err != nil {
		log.Fatalf("parse:magicNumber %v\n", err)
	}
}

func Parse2ByteValue(raw []byte, byteOrder binary.ByteOrder, container *uint16) {
	if len(raw) != 2 {
		log.Fatal("Parse2ByteValue: raw is not the right size")
	}

	buffer := bytes.NewBuffer(raw)

	err := binary.Read(buffer, byteOrder, container)

	if err != nil {
		log.Fatalf("parse:magicNumber %v\n", err)
	}
}