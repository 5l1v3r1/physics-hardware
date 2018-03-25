// +build !darwin
package command

import (
	"log"
	"fmt"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type BleList map[string]string
var DeviseList []BleDevise

type BleDevise struct{
	Name string
	UUID string
}

func addIfNotExist(Name string, UUID string) {

	found := false
	if len(DeviseList) > 0 {
		for _, value := range DeviseList {
			if value.Name == Name{
				found = true
			}
		}
		if found == false {
			DeviseList = append(DeviseList, BleDevise{Name, UUID})
			fmt.Println("Added :", DeviseList)
		}
	}
}

func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning for iBeacon Broadcasts...")
		device.Scan([]gatt.UUID{}, true)
		return
	default:
		device.StopScanning()
	}
}

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int){
	if p.Name() != "" {
		addIfNotExist(p.Name(), p.ID())
	}

	if GetStatus() == true {
		p.Device().StopScanning()
		return
	}

}

func onPeriphDisconnected(p gatt.Peripheral) {

}

func BleRecon() {
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	device.Handle(
			gatt.PeripheralDiscovered(onPeripheralDiscovered),
		)
	device.Init(onStateChanged)
}
