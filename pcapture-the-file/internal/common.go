package internal

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Parse4ByteValue(raw []byte ) uint32 {


	if len(raw) != 4 {
		log.Fatal("Parse4ByteValue: raw is not the right size")
	}

	buffer := bytes.NewBuffer(raw)

	var container uint32

	err := binary.Read(buffer, binary.LittleEndian, &container)

	if err != nil {
		log.Fatalf("parse:Parse4ByteValue %v\n", err)
	}

	return container
}

func Parse4ByteValueBig(raw []byte) uint32 {


	if len(raw) != 4 {
		log.Fatal("Parse4ByteValue: raw is not the right size")
	}

	buffer := bytes.NewBuffer(raw)

	var container uint32

	err := binary.Read(buffer, binary.BigEndian, &container)

	if err != nil {
		log.Fatalf("parse:Parse4ByteValue %v\n", err)
	}

	return container
}

func Parse2ByteValue(raw []byte) uint16 {
	if len(raw) != 2 {
		log.Fatal("Parse2ByteValue: raw is not the right size")
	}

	buffer := bytes.NewBuffer(raw)

	var container uint16

	err := binary.Read(buffer, binary.LittleEndian, &container)

	if err != nil {
		log.Fatalf("parse:Parse2ByteValue %v\n", err)
	}
	return container
}

func Parse2ByteValueBig(raw []byte) uint16 {
	if len(raw) != 2 {
		log.Fatal("Parse2ByteValue: raw is not the right size")
	}

	buffer := bytes.NewBuffer(raw)

	var container uint16

	err := binary.Read(buffer, binary.BigEndian, &container)

	if err != nil {
		log.Fatalf("parse:Parse2ByteValue %v\n", err)
	}
	return container
}

