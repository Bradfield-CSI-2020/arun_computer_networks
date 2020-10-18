package main

import (
	parse "computer_networks/pcapture-the-file/internal"
	"io/ioutil"
	"log"
)

func main() {

	data, err := ioutil.ReadFile("net.cap")

	if err != nil {
		log.Fatalf("parse:readfile %v\n", err)
	}

	parse.PcapHeader(data)
}
