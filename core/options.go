package core

import (
	"flag"
	"os"

	"github.com/graniet/physics-hardware/core/command"
)

func parseArguments() {

	ccHost := flag.String("cchost", "", "command & control host e.x: http://192.168.0.14:8080/ "+RED+"required"+RESET)
	ccGate := flag.String("ccgate", "gate/", "command & control gate e.x: gate/")
	webservice := flag.Bool("webservice", false, "load web server with root path webserver.")
	flag.Parse()

	command.CcGATE = *ccGate
	command.CcHOST = *ccHost
	userWebService = *webservice
	AskFlag = false
}

func checkFlags() bool {
	if AskFlag == false {
		return true
	} else {
		parseArguments()
		if len(os.Args) < 2 {
			flag.PrintDefaults()
		} else {
			if command.CcHOST != "" {
				return true
			} else {
				flag.PrintDefaults()
				return false
			}
		}
	}
	return false
}
