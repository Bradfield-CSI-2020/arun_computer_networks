package main

import (
	"computer_networks/parse"
	"log"
	"os"
)




func main() {

	file, fileOpenErr := os.Open("net.cap")

	var myFile parse.MyFile

	if fileOpenErr != nil {
		log.Fatalf("parse:readError %v\n", fileOpenErr)
	}

	var rawPCapHeader = make([]byte,24)

	_, pageErr := file.Read(rawPCapHeader)

	if pageErr != nil {
		log.Fatalf("something bad happened %v\n", pageErr)
	}

	globalHeader := parse.ReadPCapHeader(rawPCapHeader)

	myFile.PCapHeader = globalHeader

	//fmt.Printf("version: %d.%d\n", globalHeader.MajorVersion, globalHeader.MinorVersion)
	//fmt.Printf("timestamp offset: %d\n", globalHeader.TimestampOffset)
	//fmt.Printf("timestamp accuracy: %d\n", globalHeader.TimestampAccuracy)
	//fmt.Printf("snapshot length: %d\n", globalHeader.SnapShotLength)
	//fmt.Printf("link layer globalHeader type: %d\n", globalHeader.LinkLayerHeaderType)


	bytesRead := -1

	for bytesRead != 0 {
		var rawPacketHeader = make([]byte,16)
		read, pageErr := file.Read(rawPacketHeader)

		if pageErr != nil {
			if read != 0 {
				log.Fatalf("header page error %v\n", pageErr)
			}
			break
		}

		bytesRead = read
		packetHeader := parse.ReadPacketHeader(rawPacketHeader)

		var packetData parse.PacketData

		if packetHeader.PacketLength != packetHeader.FullPacketLength {
			log.Fatalf("packet legth mismatch %v\n", pageErr)
		}

		packetData.PacketHeader = packetHeader

		var rawPacketData = make([]byte,packetHeader.PacketLength)

		_, packetReadErr := file.Read(rawPacketData)

		if packetReadErr != nil {
			log.Fatalf("packet page error %v\n", pageErr)
		}

		packetData.RawData = rawPacketData

		if myFile.PacketDataData == nil {
			myFile.PacketDataData = []parse.PacketData{packetData}
		} else {
			myFile.PacketDataData = append(myFile.PacketDataData, packetData)
		}

		//fmt.Printf("timestamp seconds: %d\n", packetHeader.TimestampSeconds)
		//fmt.Printf("timestamp micro: %d\n", packetHeader.TimestampMicroSeconds)
		//fmt.Printf("packet length: %d\n", packetHeader.PacketLength)
		//fmt.Printf("full packet length: %d\n", packetHeader.FullPacketLength)
	}

	if len(myFile.PacketDataData) != 99 {
		log.Fatalln("total packet count mismatch")
	}

	//fmt.Printf("total packet data count: %d\n", len(myFile.PacketDataData))

}


