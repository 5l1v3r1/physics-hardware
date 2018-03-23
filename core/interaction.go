package core

import (
	"time"
	"strings"
	"fmt"

	"github.com/graniet/physics-hardware/core/command"
)

func viewWifi() []byte{

	sendCmd(apiUrlCpl, "wifi.recon on")
	time.Sleep(5000 * time.Millisecond)
	sendCmd(apiUrlCpl, "wifi.recon off")

	return sendGET(apiUrlCpl, "/api/session/wifi")
}

func checkCommand(user_command string) {

	if startedCorrect == false {
		Running()
	}

	if user_command == "physics.pcap.read" {

		command.SendEvent("command found (physics.pcap.read)","")
		if userDebug == true {
			command.ReadPCAP(true)
		} else {
			command.ReadPCAP(false)
		}
	} else if strings.Contains(user_command, "physics.pcap.filter") {

		command.SendEvent("command found (physics.pcap.filter)","")
		filterList := strings.Split(user_command, "physics.pcap.filter")
		if len(filterList) > 0 {
			filter := filterList[1]
			filter = strings.TrimSpace(filter)
			command.PacketFilter = filter
			command.SendEvent("Update filter to : " + filter, "")

			if userDebug == true{
				fmt.Println("Filter updated: " + filter)
			}
		} else {

			command.SendEvent("Filter not found or incorrect setup.", "")
			if userDebug == true{
				fmt.Println("Filter not found or incorrect setup")
			}
		}
	} else if user_command == "physics.wifi.list" {

		command.SendEvent("command found (physics.wifi.list)","")
		data := viewWifi()
		command.ReadWifi(string(data))
	} else if user_command == "physics.get.env" {

		bettercapEvent := sendGET(apiUrlCpl, "/api/session/env")
		command.ReadEnvironment(string(bettercapEvent))
		command.SendEvent("command found (physics.get.env)", "")
	}
}

func loadInteraction() {

	for {
		commandString := command.ListenCC()
		if commandString != "" {

			checkCommand(commandString)
			time.Sleep(1000 * time.Millisecond)
		}
		command.SendBattery()
	}
}
