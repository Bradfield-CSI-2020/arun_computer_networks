package main

import (
	"computer_networks/parse"
	"fmt"
)


func main() {
	header := parse.PcapHeader()

	fmt.Printf("version: %d.%d\n", header.MajorVersion, header.MinorVersion)
	fmt.Printf("timestamp offset: %d\n", header.TimestampOffset)
	fmt.Printf("timestamp accuracy: %d\n", header.TimestampAccuracy)
	fmt.Printf("snapshot length: %d\n", header.SnapShotLength)
	fmt.Printf("linke layer header type: %d\n", header.LinkLayerHeaderType)

	parse.GetPacketData()

	//fmt.Printf("timestamp seconds: %d\n", packetHeader.TimestampSeconds)
	//fmt.Printf("timestamp micro: %d\n", packetHeader.TimestampMicroSeconds)
	//fmt.Printf("packet length: %d\n", packetHeader.PacketLength)
	//fmt.Printf("full packet length: %d\n", packetHeader.FullPacketLength)
}


