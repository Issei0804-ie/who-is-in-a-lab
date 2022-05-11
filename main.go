package main

import (
	"encoding/json"
	"fmt"
	"github.com/Issei0804-ie/who-is-in-a-lab/api"
	"github.com/Issei0804-ie/who-is-in-a-lab/domain"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	device       string
	snapshot_len int32 = 1024
	promiscuous  bool  = false
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
	device = os.Args[1]
	if device == "" {
		log.Fatal("usage: ./who-is-in-a-lab [network interface] \n example: ./who-is-in-a-lab wlan0")
	}
	log.Println(fmt.Sprintf("network interface is " + device))

	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	log.Println("listen now ...")
	go api.InitAPI(&members)
	log.Println("finish listen")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacketInfo(packet, members)
	}
}

func printPacketInfo(packet gopacket.Packet, members []domain.Member) {
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
	}
}
