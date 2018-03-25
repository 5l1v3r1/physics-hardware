// +build !darwin
package command

import (
	"log"
	"fmt"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type BleList map[string]string

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
		fmt.Println("device: " + p.Name())
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

	//select {}
}
