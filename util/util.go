package util

import (
	"fmt"
	"log"
	"regexp"
)

// ValidateMacAddress mac address が正しければ true を返します
func ValidateMacAddress(macAddress string) bool {
	r, err := regexp.Compile(`^[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}$`)
	if err != nil {
		log.Println(err)
		return false
	}
	fmt.Println(macAddress)
	return r.MatchString(macAddress)
}
