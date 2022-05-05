package domain

import (
	"fmt"
	"log"
	"strconv"
	"time"
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
