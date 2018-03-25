// +build !darwin
package command

import (
	"log"
	"fmt"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type BleList map[string]string
var deviceList = BleList{}
var BleReconStop = 0


func onStateChanged(device gatt.Device, s gatt.State) {
	if GetStatus(){
		device.StopScanning()
	}
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
		deviceList[p.ID()] = p.Name()
}

func BleRecon() {
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	device.Handle(gatt.PeripheralDiscovered(onPeripheralDiscovered))
	device.Init(onStateChanged)
	select {}
}
