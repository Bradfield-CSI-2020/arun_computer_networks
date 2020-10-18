package parse

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
)

func PcapHeader(data []byte) {

	firstFour := []byte{data[0],data[1],data[2],data[3]}
	mgBuffer := bytes.NewBuffer(firstFour)

	var magicNumber uint32

	err := binary.Read(mgBuffer, binary.LittleEndian, &magicNumber)

	if err != nil {
		log.Fatalf("parse:magicNumber %v\n", err)
	}

	rep := hex.EncodeToString(firstFour)

	fmt.Printf("hex: %s\n", rep)

	fmt.Printf("magicNumber: %d\n", magicNumber)

}
//now := []byte{0xFF,0xFF,0xFF,0xFF}

//nowBuffer := bytes.NewReader(now)
//var  nowVar uint32
//binary.Read(nowBuffer,binary.BigEndian,&nowVar)
//fmt.Println(nowVar)