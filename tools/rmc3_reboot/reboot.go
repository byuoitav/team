package main

import (
	"fmt"
	"time"

	"github.com/byuoitav/team/tools/rmc3_reboot/helpers"
)

const (
	REBOOT_COMMAND  = "Reboot"
	REBOOT_RESPONSE = "Rebooting system.  Please wait...\r\n"
)

func main() {

	addresses := []string{"KMBL-680-GW1.byu.edu", "CTB-410-GW1.byu.edu",
		"HBLL-2212-GW1.byu.edu", "HBLL-2231-GW1.byu.edu",
		"EB-404-GW1.byu.edu", "HCEB-BALLRM-GW1.byu.edu",
		"JFSB-1081-GW1.byu.edu", "JRCB-267-GW1.byu.edu",
		"MCKB-150-GW1.byu.edu", "TNRB-240-GW1.byu.edu"}

	for _, ip := range addresses {
		_, err := rebootDevice(ip)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func rebootDevice(address string) (string, error) {

	connection, err := helpers.GetConnection(address, true)
	if err != nil {
		return "", err
	}

	var response string
	for ; response != REBOOT_RESPONSE; time.Sleep(2 * time.Second) {
		response, err = helpers.SendCommand(connection, REBOOT_COMMAND)
		if err != nil {
			return response, err
		}
	}

	return response, err
}
