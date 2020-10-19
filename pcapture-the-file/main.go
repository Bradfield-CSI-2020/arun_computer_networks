package main

import (
	"computer_networks/parse"
	"fmt"
	"log"
	"os"
)

func main() {

	file, fileOpenErr := os.Open("net.cap")

	var myFile parse.MyFile

	if fileOpenErr != nil {
		log.Fatalf("parse:readError %v\n", fileOpenErr)
	}

	var rawPCapHeader = make([]byte, 24)

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
		var rawPacketHeader = make([]byte, 16)
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

		var rawPacketData = make([]byte, packetHeader.PacketLength)

		_, packetReadErr := file.Read(rawPacketData)

		if packetReadErr != nil {
			log.Fatalf("packet page error %v\n", pageErr)
		}

		packetData.RawData = rawPacketData
		packetData.EtherNetFrame = parse.ReadEtherNetHeaders(rawPacketData)

		ipHeader := parse.ReadIpHeader(packetData.EtherNetFrame.IpRawPayload)

		var ipDataGram parse.IpDataGram

		rawData := packetData.EtherNetFrame.IpRawPayload[ipHeader.InternetHeaderLength:]

		ipDataGram.IpHeader = ipHeader
		ipDataGram.RawData = rawData

		if myFile.IpDataGramData == nil {
			myFile.IpDataGramData = []parse.IpDataGram{ipDataGram}
		} else {
			myFile.IpDataGramData = append(myFile.IpDataGramData, ipDataGram)
		}

		fmt.Println("")
		fmt.Printf("InternetHeaderLength: %d\n", ipHeader.InternetHeaderLength)
		fmt.Printf("Data Length %d\n", len(rawData))
		fmt.Printf("TotalLength: %d\n", ipHeader.TotalLength)

		fmt.Printf("ECN: %d\n", ipHeader.ECN)
		fmt.Printf("Protocol: %d\n", ipHeader.Protocol)
		fmt.Printf("SourceIp: %d\n", ipHeader.SourceIp)
		fmt.Printf("DestinationIp: %d\n", ipHeader.DestinationIp)
		fmt.Println("")

		// TODO: what is the proper way to decode mac addresses?
		//fmt.Printf("DestinationMac address: %s\n", hex.EncodeToString(packetData.EtherNetFrame.DestinationMac))
		//fmt.Printf("SourceMac address: %s\n", hex.EncodeToString(packetData.EtherNetFrame.SourceMac))

		// TODO: looks like requests ans responses have the mac addresses swapped
		//if "c4e984876028" != hex.EncodeToString(packetData.EtherNetFrame.DestinationMac) {
		//	log.Fatalln("DestinationMac address mismatch !")
		//}
		//
		//if "a45e60df2e1b" != hex.EncodeToString(packetData.EtherNetFrame.SourceMac) {
		//	log.Fatalln("SourceMac address mismatch !")
		//}

		//fmt.Printf("DestinationMac size: %d\n", len(packetData.EtherNetFrameRaw.DestinationMac))
		//fmt.Printf("SourceMac size: %d\n", len(packetData.EtherNetFrameRaw.SourceMac))
		//fmt.Printf("EtherType size: %d\n", len(packetData.EtherNetFrameRaw.EtherType))
		//fmt.Printf("IpRawPayload size: %d\n", len(packetData.EtherNetFrameRaw.IpRawPayload))
		//fmt.Printf("InterPacketGap size: %d\n", len(packetData.EtherNetFrameRaw.InterPacketGap))

		if myFile.PacketDataData == nil {
			myFile.PacketDataData = []parse.PacketData{packetData}
		} else {
			myFile.PacketDataData = append(myFile.PacketDataData, packetData)
		}

		totalSize := len(packetData.EtherNetFrame.DestinationMac) + len(packetData.EtherNetFrame.SourceMac) + len(packetData.EtherNetFrame.EtherType) + len(packetData.EtherNetFrame.IpRawPayload) + len(packetData.EtherNetFrame.InterPacketGap)

		if totalSize != len(packetData.RawData) {
			log.Fatalln("ethernet frame size mismatch")
		}

		//fmt.Printf("timestamp seconds: %d\n", packetHeader.TimestampSeconds)
		//fmt.Printf("timestamp micro: %d\n", packetHeader.TimestampMicroSeconds)
		//fmt.Printf("packet length: %d\n", packetHeader.PacketLength)
		//fmt.Printf("full packet length: %d\n", packetHeader.FullPacketLength)
	}

	//fmt.Printf("total packet data count: %d\n", len(myFile.PacketDataData))

}
