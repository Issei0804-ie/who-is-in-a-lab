package main

import (
	"awesomeProject/api"
	"awesomeProject/domain"
	"encoding/json"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"io/ioutil"
	"log"
	"time"
)

var (
	device       string = "enp0s31f6"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 1 * time.Second
	handle       *pcap.Handle
)

func initMembers() []domain.Member {
	var members []domain.Member
	file, err := ioutil.ReadFile("./address.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(file, &members); err != nil {
		log.Fatal(err)
	}

	for i, member := range members {
		fmt.Printf("%v \n", member)
		members[i].LastLogin = time.Time{}
	}
	return members
}

func main() {
	members := initMembers()

	log.Println("listen now ...")
	go api.InitAPI(&members)
	log.Println("finish listen")

	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process packet here
		printPacketInfo(packet, members)
	}
}

func printPacketInfo(packet gopacket.Packet, members []domain.Member) {
	// Ethernet Packetへキャスト
	// Let's see if the packet is an ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		for i, member := range members {
			for _, address := range member.Addresses {
				if ethernetPacket.SrcMAC.String() == address {
					members[i].LastLogin = time.Now()
					log.Println("packet caught:" + member.Name)
				}
			}
		}
	} else {
		fmt.Println("nil")
	}
}
