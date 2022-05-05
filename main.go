package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	device       string = "en0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 1 * time.Second
	handle       *pcap.Handle
)

type Member struct {
	Name      string
	Addresses []string
	IsLab     bool
	LastLogin time.Time
}

func (m *Member) SetIsLab(limit int) {
	now := time.Now()
	log.Printf("name = %v, loginTime = %v, now is %v \n", m.Name, m.LastLogin, now)
	log.Printf("sub is " + strconv.Itoa(int(now.Sub(m.LastLogin).Minutes())) + "\n")
	fmt.Println("sub" + strconv.Itoa(int(now.Sub(m.LastLogin).Minutes())))
	if now.Sub(m.LastLogin).Minutes() <= float64(limit) {
		m.IsLab = true
		fmt.Println(m)
	} else {
		m.IsLab = false
	}
}

func initMembers() []Member {
	var members []Member
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		limit := 30
		for i := 0; i < len(members); i++ {
			members[i].SetIsLab(limit)
		}
		t := template.Must(template.ParseFiles("index.html"))
		t.Execute(w, members)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Write([]byte("not allow this method"))
			return
		}

		newMember := Member{}
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &newMember)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Printf("%v, \n", newMember)
		if err != nil || newMember.Name == "" || newMember.Addresses == nil {

			log.Println(err)
			log.Println(string(body))
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = "invalid body"
			jsonResp, _ := json.Marshal(resp)
			w.Write(jsonResp)
			return
		}

		// mac address が既に登録されていないか確認
		for _, member := range members {
			for _, newMemberAddress := range newMember.Addresses {
				for _, address := range member.Addresses {
					if address == newMemberAddress {
						w.WriteHeader(http.StatusBadRequest)
						w.Header().Set("Content-Type", "application/json")
						resp := make(map[string]string)
						resp["message"] = "this mac address is already used. if you want to remove stored mac address, you need to ask issei."
						jsonResp, _ := json.Marshal(resp)
						w.Write(jsonResp)
						return
					}
				}
			}
		}

		didAdd := false
		// 同じ名前なら mac address を追加
		for i, member := range members {
			if newMember.Name == member.Name {
				members[i].Addresses = append(member.Addresses, newMember.Addresses...)
				didAdd = true
			}
		}

		if !didAdd {
			members = append(members, newMember)
		}
		jsonMembers, err := json.Marshal(members)
		if err != nil {
			log.Println(err)
			return
		}
		file, err := os.Create("./address.json")
		defer file.Close()
		if err != nil {
			log.Println(err)
			return
		}

		file.Write(jsonMembers)
		file.Sync()

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "ok"
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
		return

	})
	log.Println("listen now ...")

	go http.ListenAndServe(":80", nil)

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

func printPacketInfo(packet gopacket.Packet, members []Member) {
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
