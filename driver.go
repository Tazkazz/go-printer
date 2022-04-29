package go_printer

import (
	"errors"
	"fmt"
	"github.com/google/gousb"
	"log"
	"time"
)

const deviceClass = gousb.ClassPrinter

// TODO: allow searching for printer devices without known id
// TODO: automatically search for printer after disconnection

func NewVoidPrinter() (printer *UsbPrinter) {
	printer = &UsbPrinter{
		write: func(data []byte) (wrote int, err error) {
			return 0, nil
		},
		state: &UsbPrinterState{},
		Done:  func() {},
	}
	printer.initialize()
	return
}

func NewUsbPrinter(vendorId, productId uint16) (printer *UsbPrinter, err error) {
	log.Printf("Searching for a printer with vendorId: 0x%x, productId: 0x%x", vendorId, productId)
	endp, endpDone, err := openOutEndpoint(vendorId, productId)
	if err != nil {
		return
	}
	log.Printf("Printer was found")

	printer = &UsbPrinter{
		write: endp.Write,
		state: &UsbPrinterState{},
		Done:  endpDone,
	}

	printer.initialize()

	return
}

func openOutEndpoint(vendorId, productId uint16) (endpoint *gousb.OutEndpoint, done func(), err error) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	dev, err := ctx.OpenDeviceWithVIDPID(gousb.ID(vendorId), gousb.ID(productId))
	if err != nil {
		return
	}

	intf, intfDone, err := dev.DefaultInterface()
	if err != nil {
		return
	}

	if intf.Setting.Class != deviceClass {
		defer intfDone()
		err = errors.New(fmt.Sprintf("device class %s is not a %s", intf.Setting.Class, deviceClass))
		return
	}

	description := findOutEndpointDesc(intf)
	if description == nil {
		defer intfDone()
		err = errors.New("device doesn't have any OUT endpoints")
		return
	}

	endp, err := intf.OutEndpoint(int(description.Address))
	if err != nil {
		defer intfDone()
		return
	}

	descIn := findInEndpointDesc(intf)
	if descIn == nil {
		defer intfDone()
		err = errors.New("device doesn't have any IN endpoints")
		return
	}

	endIn, err := intf.InEndpoint(int(description.Address))
	if err != nil {
		defer intfDone()
		return
	}

	go func() {
		b := make([]byte, 1)
		log.Println("Starting to read from the printer...")
		for {
			//log.Println("Reading from the printer")
			s, _ := endIn.NewStream(1, 1)
			//i, _ := endIn.Read(b)
			i, _ := s.Read(b)
			time.Sleep(5000000)
			//if i == 0 {
			//	continue
			//}
			//log.Println("Read!")
			log.Println(i, b)
		}
	}()

	endpoint = endp
	done = intfDone
	return
}

func findOutEndpointDesc(intf *gousb.Interface) *gousb.EndpointDesc {
	for _, description := range intf.Setting.Endpoints {
		if description.Direction == gousb.EndpointDirectionOut {
			return &description
		}
	}
	return nil
}

func findInEndpointDesc(intf *gousb.Interface) *gousb.EndpointDesc {
	for _, description := range intf.Setting.Endpoints {
		if description.Direction == gousb.EndpointDirectionIn {
			return &description
		}
	}
	return nil
}
